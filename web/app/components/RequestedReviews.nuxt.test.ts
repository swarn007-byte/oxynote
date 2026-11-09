import { mountSuspended } from "@nuxt/test-utils/runtime"
import { enableAutoUnmount } from "@vue/test-utils"
import { afterEach, beforeEach, describe, it } from "vitest"
import { createError } from "h3"
import {
	clearQueryCache,
	disposeMockEndpoints,
	mockDeferredEndpoint,
	mockEndpoint,
} from "~/composables/api/test-helpers"
import { seedAuthOrganization, settleMutations, t } from "./test-helpers"
import RequestedReviews from "./RequestedReviews.vue"

enableAutoUnmount(afterEach)

describe("RequestedReviews", { concurrent: false }, () => {
	// each mount uses the shared app's query cache and endpoint registry.
	beforeEach(() => {
		clearQueryCache()
		seedAuthOrganization({ slug: "Acme" })
	})
	afterEach(disposeMockEndpoints)

	it("links each pending review to its exact branch", async ({ expect }) => {
		const calls = mockEndpoint("GET", "/api/reviews", () => [
			{
				documentId: "doc1".padEnd(20, "0"),
				branchId: "draft1",
				documentName: "Runbook",
				branchName: "Update",
				icon: "lucide:file",
			},
		])
		const wrapper = await mountSuspended(RequestedReviews)
		await settleMutations()

		const link = wrapper.get("a")
		expect(link.text()).toContain("Runbook")
		expect(link.text()).toContain("Update")
		expect(link.attributes("href")).toBe(
			`/Acme/Runbook-${"doc1".padEnd(20, "0")}?branch=draft1`,
		)
		expect(calls).toHaveLength(1)
	})

	it("shows an empty state when no reviews are pending", async ({ expect }) => {
		const calls = mockEndpoint("GET", "/api/reviews", () => [])
		const wrapper = await mountSuspended(RequestedReviews)
		await settleMutations()

		expect(wrapper.text()).toContain(t("sidebar.reviews.empty"))
		expect(wrapper.find("a").exists()).toBe(false)
		expect(calls).toHaveLength(1)
	})

	it("shows loading until the request finishes", async ({ expect }) => {
		const endpoint = mockDeferredEndpoint("GET", "/api/reviews")
		const wrapper = await mountSuspended(RequestedReviews)
		await endpoint.reached

		expect(wrapper.get('[role="status"]').text()).toBe(
			t("sidebar.reviews.loading"),
		)
		endpoint.resolve([])
		await settleMutations()
		expect(wrapper.text()).toContain(t("sidebar.reviews.empty"))
		expect(endpoint.calls).toHaveLength(1)
	})

	it("reports a failure and lets the user retry", async ({ expect }) => {
		let failed = true
		const calls = mockEndpoint("GET", "/api/reviews", () => {
			if (failed) {
				throw createError({ statusCode: 400, message: "unavailable" })
			}

			return []
		})
		const wrapper = await mountSuspended(RequestedReviews)
		await settleMutations()
		expect(wrapper.get('[role="alert"]').text()).toBe(
			t("sidebar.reviews.error"),
		)
		failed = false

		await wrapper.get("button").trigger("click")
		await settleMutations()

		expect(wrapper.text()).toContain(t("sidebar.reviews.empty"))
		expect(wrapper.find('[role="alert"]').exists()).toBe(false)
		expect(calls).toHaveLength(2)
	})
})
