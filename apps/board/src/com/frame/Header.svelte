<div class="header-root">
  <div class="header-title">
    <div class="btn-menu">
      <div class="menu icon icon-menu">
        <IconMenu/>
      </div>
    </div>
    <div class="title">Sample document</div>
  </div>
  <div class="header-middle">

  </div>
  <div class="header-right">
    <div class="btn-button btn-primary">Publish</div>
  </div>
</div>

<script lang="ts">
  import IconMenu from '../../assets/IconMenu.svelte';
</script>

<style lang="scss">
  @use '../../styles/mixins' as m;

  .header-root {
    @include m.flex-row;
    flex: 1 0 auto;
    align-items: center;
    justify-content: center;
  }

  .header-title {
    @include m.flex-row;
    flex: 1 0 auto;
    height: 100%;

    .btn-menu {
      @include m.flex-row;
      align-items: center;
      height: 100%;
      cursor: pointer;

      &:hover {

      }
    }

    .title {
      line-height: 40px;
    }

    .icon-menu {
      height: 20px;
      width: 20px;
      margin: 0 10px;
    }
  }

  .header-middle {
    flex: 1 0 auto;
  }

  .header-right {
    flex: 0 0 auto;
    margin-right: 10px;
  }
</style>
