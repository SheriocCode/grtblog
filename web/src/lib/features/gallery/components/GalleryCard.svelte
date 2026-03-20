<script lang="ts">
	import type { GallerySummary } from '$lib/features/gallery/types';

	interface Props {
		item: GallerySummary;
		onPreview?: (src: string, rect: DOMRect | null, alt: string) => void;
	}

	let { item, onPreview }: Props = $props();

	const createdAt = $derived.by(() => new Date(item.createdAt));
	const dayLabel = $derived.by(() => String(createdAt.getDate()).padStart(2, '0'));
	const timeLabel = $derived.by(
		() => `${String(createdAt.getHours()).padStart(2, '0')}:${String(createdAt.getMinutes()).padStart(2, '0')}`
	);
	const monthLabel = $derived.by(
		() => `${String(createdAt.getMonth() + 1).padStart(2, '0')}月`
	);
	const lines = $derived.by(() =>
		item.content
			.split(/\r?\n/)
			.map((line) => line.trim())
			.filter(Boolean)
	);

	const gridClass = $derived.by(() => {
		const count = item.images?.length ?? 0;
		if (count <= 1) return 'grid-cols-1';
		if (count === 2) return 'grid-cols-2';
		if (count === 4) return 'grid-cols-2';
		return 'grid-cols-3';
	});

	function handlePreview(event: MouseEvent, src: string, index: number) {
		event.stopPropagation();
		const rect = event.currentTarget instanceof HTMLElement ? event.currentTarget.getBoundingClientRect() : null;
		onPreview?.(src, rect, `Gallery ${item.id} image ${index + 1}`);
	}
</script>


<article class="relative border-b border-dashed border-ink-200 py-6 last:border-b-0 dark:border-ink-700">
	<div class="relative flex gap-4">
		<div class="flex w-16 shrink-0 flex-col items-center px-2 py-2 text-center">
			<div class="font-mono text-[10px] uppercase tracking-[0.28em] text-jade-600 dark:text-jade-400">{monthLabel}</div>
			<div class="mt-1 font-serif text-3xl font-bold text-ink-900 dark:text-ink-100">{dayLabel}</div>
			<div class="mt-1 font-mono text-[10px] tracking-[0.18em] text-ink-400 dark:text-ink-500">{timeLabel}</div>
			{#if item.isTop}
				<div class="mt-3 rounded-full border border-jade-500/30 px-2 py-1 font-mono text-[9px] uppercase tracking-[0.22em] text-jade-700 dark:border-jade-500/20 dark:text-jade-300">Top</div>
			{/if}
		</div>

		<div class="min-w-0 flex-1">
			<div class="mb-4 flex items-center justify-between gap-3">
				<div class="font-mono text-[10px] uppercase tracking-[0.26em] text-ink-400 dark:text-ink-500">Daily Frame #{item.id}</div>
				<div class="text-[11px] text-ink-400 dark:text-ink-500">{item.imageCount} 张图</div>
			</div>

			<div class="space-y-2">
				{#each lines as line, lineIndex (lineIndex)}
					<p class="font-serif text-[15px] leading-8 text-ink-700 dark:text-ink-300">{line}</p>
				{/each}
			</div>

			{#if item.images && item.images.length > 0}
				<div class={`mt-5 grid ${gridClass} gap-2 md:gap-3`}>
					{#each item.images as src, index (src + index)}
						<button
							type="button"
							class="group/image relative overflow-hidden rounded-2xl border border-ink-200 bg-transparent transition-colors duration-300 hover:border-ink-300 dark:border-ink-700 dark:bg-transparent dark:hover:border-ink-600"
							onclick={(event) => handlePreview(event, src, index)}
						>
							<img
								src={src}
								alt={`Gallery ${item.id} image ${index + 1}`}
								class={`h-full w-full object-cover transition-transform duration-500 group-hover/image:scale-[1.03] ${(item.images?.length ?? 0) === 1 ? 'aspect-[16/9]' : 'aspect-square'}`}
								loading="lazy"
							/>						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</article>

<style lang="postcss">
	@reference "$routes/layout.css";
</style>