export interface RequestedReview {
	documentId: string
	branchId: string
	documentName: string
	branchName: string
	icon: string
}

export default function () {
	const { $coreAPIClient } = useNuxtApp()
	const { fetchAuthSession } = useAuthSession()

	const fetchRequestedReviews = useQuery({
		key: () => [
			"requested-reviews",
			fetchAuthSession.data.value?.data?.session.activeOrganizationId ?? "",
			fetchAuthSession.data.value?.data?.session.userId ?? "",
		],
		query: () =>
			$coreAPIClient<RequestedReview[]>("/api/reviews", { method: "GET" }),
		refetchOnMount: true,
		refetchOnWindowFocus: true,
		staleTime: 0,
	})

	return { fetchRequestedReviews }
}
