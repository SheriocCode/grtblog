<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import { galleryListCtx } from '$lib/features/gallery/context';
	import GalleryTimeline from '$lib/features/gallery/components/GalleryTimeline.svelte';
	import Pagination from '$lib/ui/primitives/pagination/Pagination.svelte';
	import PageHeader from '$lib/ui/common/PageHeader.svelte';
	import ImagePreview from '$lib/ui/markdown/ImagePreview.svelte';
	import { resolvePath } from '$lib/shared/utils/resolve-path';

	let { data } = $props<{ data: PageData }>();

	galleryListCtx.mountModelData(() => data.galleries);

	const items = galleryListCtx.selectModelData((d) => d?.items || []);
	const total = galleryListCtx.selectModelData((d) => d?.total ?? 0);
	const page = galleryListCtx.selectModelData((d) => d?.page ?? 1);
	const size = galleryListCtx.selectModelData((d) => d?.size ?? 12);

	const totalPages = $derived($size > 0 ? Math.max(1, Math.ceil($total / $size)) : 1);

	const groups = $derived.by(() => {
		const grouped = new Map<string, { key: string; title: string; subtitle: string; items: typeof $items }>();
		for (const item of $items) {
			const date = new Date(item.createdAt);
			const key = `${date.getFullYear()}-${date.getMonth() + 1}`;
			const existing = grouped.get(key);
			const group = existing ?? {
				key,
				title: `${date.getFullYear()}年 ${String(date.getMonth() + 1).padStart(2, '0')}月`,
				subtitle: date.getFullYear() === new Date().getFullYear() ? 'THIS SEASON' : 'ARCHIVE FRAME',
				items: []
			};
			group.items.push(item);
			grouped.set(key, group);
		}
		return Array.from(grouped.values());
	});

	let preview = $state<{ src: string; rect: DOMRect | null; alt: string } | null>(null);

	function handlePreview(src: string, rect: DOMRect | null, alt: string) {
		preview = { src, rect, alt };
	}

	function closePreview() {
		preview = null;
	}

	const onPageChange = (nextPage: number) => {
		const safePage = Number.isFinite(nextPage) && nextPage > 1 ? nextPage : 1;
		const query = safePage === 1 ? '/gallery/' : `/gallery/?page=${safePage}`;
		goto(resolvePath(query));
	};
</script>

<div class="relative mx-auto w-full max-w-6xl px-4 py-16 md:px-6">

	<div class="relative">
		<PageHeader
			title="日常 Gallery"
			tag="Gallery"
			subtitle="把普通的一天，排成有呼吸感的时间流"
			description="向内生长，莫向外求。"
		/>

		{#if $items.length > 0}
			<GalleryTimeline {groups} onPreview={handlePreview} />

			{#if totalPages > 1}
				<div class="mt-12 flex justify-center">
					<Pagination current={$page} total={totalPages} {onPageChange} />
				</div>
			{/if}
		{:else}
			<div class="flex min-h-[320px] flex-col items-center justify-center rounded-3xl border border-dashed border-ink-200 bg-white/70 px-6 text-center text-ink-400 dark:border-ink-800 dark:bg-ink-950/50 dark:text-ink-500">
				<div class="font-mono text-[11px] uppercase tracking-[0.35em]">Gallery</div>
				<p class="mt-4 font-serif text-lg">还没有公开的日常。</p>
			</div>
		{/if}
	</div>
</div>

{#if preview}
	<ImagePreview src={preview.src} alt={preview.alt} originRect={preview.rect} onClose={closePreview} />
{/if}

<style lang="postcss">
	@reference "$routes/layout.css";
</style>