<div class="editor-root">
  <div class="editor-canvas-wrap">
    <div class="editor-canvas">
      <div class="bg-canvas"></div>
      <CanvasNodes/>
    </div>
  </div>
  <div class="editor-panel">
    <PanelTabs/>
  </div>
</div>

<script lang="ts">
  import CanvasNodes from './CanvasNodes.svelte';
  import PanelTabs from './PanelTabs.svelte';
</script>

<style lang="scss">
  @use '../../styles/mixins' as m;

  .editor-root {
    @include m.flex-row;
    flex: 1 0 0;
  }

  .editor-canvas {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .bg-canvas {
    position: absolute;
    width: 10000px;
    height: 10000px;
    top: -5000px;
    left: -5000px;
    background-size: 50px 50px;
    background-image: linear-gradient(90deg, #e0e0e0 1px, transparent 0px), linear-gradient(#e0e0e0 1px, transparent 0px);
    background-color: var(--faint-gray);
    z-index: -1;
  }

  .editor-canvas-wrap {
    flex: 1 0 0;
  }

  .editor-panel {
    @include m.flex-col;
    flex: 0 0 300px;
    background: var(--white);
    border-left: solid 1px var(--border-gray);
  }
</style>
