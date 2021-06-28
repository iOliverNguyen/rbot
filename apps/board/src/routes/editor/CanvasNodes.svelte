<div class="board-canvas">
  {#each $nodes.nodes as node}
    <CanvasNode node={node}/>
  {/each}
</div>

<script lang="ts">
  import CanvasNode from './CanvasNode.svelte';
  import {nodes} from '../../stores/editor.js';
</script>

<style lang="scss">
  .board-canvas {
    position: relative;
    width: 100%;
    height: 100%;
    left: 50%;
    top: 50%;
  }


</style>
