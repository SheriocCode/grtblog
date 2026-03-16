<script lang="ts">
	import PostList from '$lib/features/post/components/PostList.svelte';
	import { postListCtx } from '$lib/features/post/context';
	import PageHeader from '$lib/ui/common/PageHeader.svelte';

	let { data } = $props();

	postListCtx.mountModelData(() => data);
</script>

<div class="max-w-5xl mx-auto py-10">
	<PageHeader
		title="Coding 笔记"
		tag="Archive"
		subtitle="指尖敲下的星辰"
		description="以时间为轴，收纳每一次逻辑的编织、思路的沉淀，留存成长的轨迹。"
	/>
	<PostList />
</div>
