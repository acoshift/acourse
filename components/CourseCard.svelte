<script>
	export let title
	export let type
	export let link
	export let desc
	export let image
	export let startDate
	export let price
	export let showDiscount
	export let discount
</script>

<style lang="scss">
	import ''

	.acourse-card {
		width: 100%;
		border-radius: 5px;
		background-color: white;
		box-shadow: 0 15px 35px rgba(50,50,93,.03), 0 5px 15px rgba(0,0,0,.06);
		margin-bottom: 16px;
		overflow: hidden;
			&.-hover-rise {
				 transition: 0.16s all ease;
				 cursor: pointer;
			&:hover {
				 box-shadow: 0 15px 35px rgba(50,50,93,.07), 0 5px 15px rgba(0,0,0,.1);
				 transform: translate(0, -4px);
			 }
		}
	}

	.course-card {
		.course-cover {
			flex-basis: 160px;
			height: 160px;
			min-height: 0;
			object-fit: cover;
		}

		.acourse-label.Live {
			background-color: #ff4e4e;
		}

		.acourse-label.eBook {
			background-color: #55da65;
		}

		.acourse-label.Video {
			background-color: #594eff;
		}

		.live-date {
			margin-left: 8px;
			color: gray;
			display: inline-block;
		}

		.course-price .discount {
			text-decoration: line-through;
			font-weight: 400;
			font-size: 13px;
		}

		.placeholder {
			background-image: url(/-/placeholder-img.svg);
			background-size: cover;
		}
	}
</style>

<svelte:options tag="course-card" />
<h1>Hello {name}!</h1>
<a href={link} class="course-card acourse-card -hover-rise _flex-column _color-dark">
<img class="course-cover placeholder" src={image} width="100%">
<div class="course-detail acourse-segment _flex-column _flex-span _main-start">
	<h4>{title}</h4>

	<div class="acourse-block">
		<div class="{`acourse-label _font-bold ${type}`}">{type}</div>
		{#if startDate}
			<div class="live-date _font-size-small">เริ่มเรียน {startDate}</div>
		{/if}
	</div>

	<div class="acourse-block-big _flex-span _font-sub _font-size-normal">{desc}</div>
	<div class="_flex-row _main-space-between">
		<div></div>
		<div class="course-price _flex-row _cross-end">
			{#if price < 0}
			<p class="_font-bold _no-margin">ฟรี</p>
			{:else}
			<p class="_no-margin"
				class:discount={showDiscount}
				class:_color-positive={showDiscount}
				class:_font-bold={!showDiscount}>
				฿ {price}
				{#if showDiscount}
				<p class="_font-bold _no-margin">&nbsp;฿ {discount}</p>
				{/if}
				</p>
			{/if}
		</div>
	</div>
</div>
</a>
