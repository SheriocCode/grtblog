import { getGalleryList } from '$lib/features/gallery/api';
import { trackISRDeps } from '$lib/server/isr-deps';
import type { PageServerLoad } from './$types';

const TRACKED_GALLERY_LIST_PAGES = 3;

export const load: PageServerLoad = async (event) => {
	const { fetch, url } = event;
	const page = Number(url.searchParams.get('page')) || 1;
	const pageSize = Number(url.searchParams.get('pageSize')) || 12;
	if (page <= TRACKED_GALLERY_LIST_PAGES) {
		trackISRDeps(event, `gallery:list:page:${page}`);
	}
	trackISRDeps(event, 'gallery:recent');

	const galleries = await getGalleryList(fetch, { page, pageSize });
	return { galleries };
};