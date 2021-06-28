<div class="list-wrap panel-insert">
  <div class="list">
    <ItemComponent type="webhook"/>
    <ItemComponent type="message"/>
    <ItemComponent type="timer"/>
    <ItemComponent type="send"/>
    <ItemComponent type="sample"/>
  </div>
</div>

<script lang="ts">
  import ItemComponent from './InsertNode.svelte';
</script>

<style lang="scss">
  .list-wrap {
    flex: 1 0 0;
    overflow: hidden auto;
  }

  .list {

  }
</style>
