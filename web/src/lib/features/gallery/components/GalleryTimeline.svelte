<script lang="ts">
	import GalleryCard from '$lib/features/gallery/components/GalleryCard.svelte';
	import StaggerList from '$lib/ui/animation/StaggerList.svelte';
	import type { GallerySummary } from '$lib/features/gallery/types';

	interface GalleryGroup {
		key: string;
		title: string;
		subtitle: string;
		items: GallerySummary[];
	}

	interface Props {
		groups: GalleryGroup[];
		onPreview?: (src: string, rect: DOMRect | null, alt: string) => void;
	}

	let { groups, onPreview }: Props = $props();
</script>

<div class="space-y-14">
	{#each groups as group (group.key)}
		<section class="relative">
			<div class="mb-6 flex items-end justify-between gap-4 border-b border-ink-200/70 pb-4 dark:border-ink-800/80">
				<div>
					<div class="font-mono text-[11px] uppercase tracking-[0.32em] text-jade-600 dark:text-jade-400">{group.subtitle}</div>
					<h2 class="mt-2 font-serif text-3xl font-bold text-ink-950 dark:text-ink-50">{group.title}</h2>
				</div>
				<div class="rounded-full border border-ink-200/70 px-3 py-1 font-mono text-[10px] uppercase tracking-[0.22em] text-ink-400 dark:border-ink-800/80">{group.items.length} 帧</div>
			</div>

			<StaggerList class="space-y-5" staggerDelay={70} duration={500} y={14} key={`gallery-${group.key}`}>
				{#each group.items as item (item.id)}
					<GalleryCard {item} {onPreview} />
				{/each}
			</StaggerList>
		</section>
	{/each}
</div>

<style lang="postcss">
	@reference "$routes/layout.css";
</style>