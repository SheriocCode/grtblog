export type GallerySummary = {
	id: number;
	content: string;
	contentHash: string;
	images?: string[];
	imageCount: number;
	isPublished: boolean;
	isTop: boolean;
	createdAt: string;
	updatedAt: string;
};

export type GalleryDetail = GallerySummary & {
	authorId: number;
	extInfo?: Record<string, unknown> | null;
};

export type GalleryListResponse = {
	items: GallerySummary[];
	total: number;
	page: number;
	size: number;
};