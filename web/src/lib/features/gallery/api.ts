import { getApi } from '$lib/shared/clients/api';
import type { GalleryListResponse } from '$lib/features/gallery/types';

type GalleryListOptions = {
	page?: number;
	pageSize?: number;
};

export const getGalleryList = async (
	fetcher?: typeof fetch,
	{ page = 1, pageSize = 12 }: GalleryListOptions = {}
): Promise<GalleryListResponse> => {
	const api = getApi(fetcher);
	const query = new URLSearchParams({
		page: String(page),
		pageSize: String(pageSize)
	});
	const result = await api<GalleryListResponse>(`/galleries?${query.toString()}`);
	return result ?? { items: [], total: 0, page, size: pageSize };
};

export const getRecentGalleries = async (fetcher?: typeof fetch): Promise<GalleryListResponse> => {
	const api = getApi(fetcher);
	const result = await api<GalleryListResponse>('/public/galleries/recent');
	return result ?? { items: [], total: 0, page: 1, size: 6 };
};