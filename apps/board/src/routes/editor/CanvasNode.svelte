<div class="canvas-node {getClass(decl?.kind)}"
     style="left: {node?.x}px; top: {node?.y}px"
     on:click={()=>clickNode()}>
  <div class="icon">
    <svelte:component this={decl?.icon}/>
  </div>
  <div class="title">{decl?.title}</div>
</div>

<script lang="ts">
  import type {Node} from '../../core/nodes.js';
  import {nodeByType} from '../../core/nodes.js';
  import {activePanelTab} from '../../stores/editor.js';

  export let node: Node;

  let decl = nodeByType(node?.type);

  const getClass = (kind: string) => 'node-' + kind;

  function clickNode() {
    console.log('click node', node);
    $activePanelTab = 'details'
  }
</script>

<style lang="scss">
  .canvas-node {
    position: absolute;
    width: 100px;
    height: 100px;
    background: var(--near-white);
    border-radius: 8px;
    border: dashed 2px transparent;
    cursor: pointer;
    z-index: 1;

    &.node-trigger {
      border: dashed 2px var(--dark-green);
    }

    &.node-action {
      border: dashed 2px var(--dark-orange);
    }
  }

  .title {
    position: absolute;
    bottom: 0;
    width: 96px;
    height: 16px;
    text-align: center;
    font-size: 12px;
  }

  .icon {
    position: absolute;
    left: 23px;
    top: 20px;
    width: 50px;
    height: 50px;
  }
</style>
