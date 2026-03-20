import { createModelDataContext } from 'svatoms';
import type { GalleryListResponse } from '$lib/features/gallery/types';

export const galleryListCtx = createModelDataContext<GalleryListResponse>({
	name: 'galleryListCtx',
	initial: null
});