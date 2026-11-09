import { afterEach, beforeEach, describe, it } from "vitest"
import { createError } from "h3"
import {
	clearQueryCache,
	disposeMockEndpoints,
	mockEndpoint,
	runInApp,
} from "./test-helpers"
import useReviewAPI from "./useReviewAPI"

describe("useReviewAPI", { concurrent: false }, () => {
	// the endpoint registry and query cache belong to the shared app.
	beforeEach(clearQueryCache)
	afterEach(disposeMockEndpoints)

	it("fetches requested reviews without caller-controlled identity", async ({
		expect,
	}) => {
		const reviews = [
			{
				documentId: "doc",
				branchId: "draft",
				documentName: "Runbook",
				branchName: "draft",
				icon: "lucide:file",
			},
		]
		const calls = mockEndpoint("GET", "/api/reviews", () => reviews)
		const api = runInApp(() => useReviewAPI())

		const result = await api.fetchRequestedReviews.refresh()

		expect(result.data).toEqual(reviews)
		expect(calls).toHaveLength(1)
		expect(calls[0]?.query).toEqual({})
	})

	it("exposes a failed request for the page to report", async ({ expect }) => {
		const calls = mockEndpoint("GET", "/api/reviews", () => {
			throw createError({ statusCode: 400, message: "unavailable" })
		})
		const api = runInApp(() => useReviewAPI())

		const result = await api.fetchRequestedReviews.refresh()

		expect(result.status).toBe("error")
		expect(result.error).toBeTruthy()
		expect(calls).toHaveLength(1)
	})
})
