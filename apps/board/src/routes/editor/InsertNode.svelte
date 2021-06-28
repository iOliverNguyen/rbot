<div class="list-item insert-component"
     on:click={_insertNode}>
  <div class="icon">
    <svelte:component this={icon}/>
  </div>
  <div class="content">
    <div class="title">{title}</div>
    <div class="desc">{desc}</div>
  </div>
</div>

<script lang="ts">
  import type {Node, NodeType} from '../../core/nodes.js';
  import {insertNode, newNode, nodeByType, nodeTypeInvalid} from '../../core/nodes.js';
  import {nodes} from '../../stores/editor.js';

  export let type: NodeType;

  let decl = nodeByType(type);
  let {icon, title, desc} = decl || nodeTypeInvalid;

  function _insertNode() {
    const node: Node = newNode(type);

    insertNode($nodes, node);
    $nodes = $nodes; // trigger rendering
  }
</script>

<style lang="scss">
  @use '../../styles/mixins' as m;

  .insert-component {
    @include m.flex-row;
    @include m.flex-center-col;

    height: 60px;
    border-bottom: solid 1px var(--washed-gray);
    cursor: pointer;
  }

  .icon {
    width: 40px;
    height: 40px;
    margin-left: 10px;
  }

  .content {
    flex: 1 0 0;
    padding: 15px 10px;
  }

  .title {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 6px;
  }

  .desc {
    font-size: 12px;
    color: var(--dark-gray);
  }

  .insert-component:hover {
    .title {
      color: var(--dark-blue);
    }
  }
</style>
