<script lang="ts" setup>
import useReviewAPI, {
	type RequestedReview,
} from "~/composables/api/useReviewAPI"

const { t } = useI18n({ useScope: "global" })
const { fetchRequestedReviews } = useReviewAPI()
const { fetchOrganization } = useAuthSession()

function reviewHref(review: RequestedReview) {
	const organization = fetchOrganization.data.value?.data?.slug ?? ""

	return `/${createNameSlug(organization)}/${createNameSlugWithId(review.documentName, review.documentId)}?branch=${encodeURIComponent(review.branchId)}`
}
</script>

<template>
	<section
		class="mx-auto max-w-4xl px-4 py-8 sm:px-6"
		aria-labelledby="reviews-title"
	>
		<div
			class="mb-6 flex flex-col items-start justify-between gap-4 sm:flex-row"
		>
			<div>
				<h1 id="reviews-title" class="text-2xl font-semibold">
					{{ t("sidebar.reviews.title") }}
				</h1>
				<p class="mt-2 text-sm text-muted-foreground">
					{{ t("sidebar.reviews.description") }}
				</p>
			</div>
			<ShadcnUiButton
				variant="outline"
				size="sm"
				:disabled="fetchRequestedReviews.isLoading.value"
				@click="fetchRequestedReviews.refetch()"
			>
				{{ t("sidebar.reviews.refresh") }}
			</ShadcnUiButton>
		</div>
		<p v-if="fetchRequestedReviews.status.value === 'error'" role="alert">
			{{ t("sidebar.reviews.error") }}
		</p>
		<p v-else-if="fetchRequestedReviews.isPending.value" role="status">
			{{ t("sidebar.reviews.loading") }}
		</p>
		<p
			v-else-if="!fetchRequestedReviews.data.value?.length"
			class="rounded-lg border border-dashed p-8 text-center text-muted-foreground"
		>
			{{ t("sidebar.reviews.empty") }}
		</p>
		<ul v-else class="divide-y rounded-lg border">
			<li
				v-for="review in fetchRequestedReviews.data.value"
				:key="review.branchId"
			>
				<NuxtLink
					:to="reviewHref(review)"
					class="flex items-center justify-between gap-4 rounded-lg p-4 hover:bg-accent focus-visible:ring-2 focus-visible:ring-ring"
				>
					<span class="min-w-0 flex-1">
						<span class="block truncate font-medium">
							{{ review.documentName }}
						</span>
						<span class="block truncate text-sm text-muted-foreground">
							{{ review.branchName }}
						</span>
					</span>
					<Icon name="lucide:arrow-right" class="shrink-0" aria-hidden="true" />
				</NuxtLink>
			</li>
		</ul>
	</section>
</template>
