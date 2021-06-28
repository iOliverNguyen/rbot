<div class="frame-root">
  <div class="frame-header">
    <Header/>
  </div>
  <div class="frame-body">
    <svelte:component this={component}/>
  </div>
  <div class="frame-footer"></div>
</div>

<script lang="ts">
  import Header from './Header.svelte';

  export let component;
</script>

<style lang="scss">
  @use '../../styles/mixins' as m;

  .frame-root {
    @include m.flex-col;
    flex: 1 0 0;
  }

  .frame-header {
    @include m.flex-row;
    flex: 0 0 40px;
    border-bottom: solid 1px var(--border-gray);
    background: var(--white);
    box-shadow: 0 0 2px var(--light-gray);
  }

  .frame-body {
    @include m.flex-row;
    flex: 1 0 auto;
  }

  .frame-footer {

  }
</style>
