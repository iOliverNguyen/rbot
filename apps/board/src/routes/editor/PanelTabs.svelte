<div class="panel-tabs">
  <div class="panel-tab"
       class:tab-active={$activePanelTab === 'insert'}
       on:click={()=>selectTab('insert')}>
    <div class="icon icon-add">
      <IconAddCircle/>
    </div>
    <div class="text">Insert</div>
  </div>
  <div class="panel-tab"
       class:tab-active={$activePanelTab === 'details'}
       on:click={()=>selectTab('details')}>
    <div class="icon icon-layer">
      <IconLayer/>
    </div>
    <div class="text">Details</div>
  </div>
</div>
<div class="panel-tabs-body">
  <svelte:component this={tabComponents[$activePanelTab]}/>
</div>

<script lang="ts">
  import IconAddCircle from '../../assets/IconAddCircle.svelte';
  import IconLayer from '../../assets/IconLayer.svelte';
  import {activePanelTab} from '../../stores/editor.js';
  import PanelDetails from './PanelDetails.svelte';
  import PanelInsert from './PanelInsert.svelte';

  const tabComponents = {
    insert: PanelInsert,
    details: PanelDetails,
  };

  const selectTab = (name: string) => {
    $activePanelTab = name;
  };
</script>

<style lang="scss">
  @use '../../styles/mixins' as m;

  .panel-tabs {
    @include m.flex-row;
    align-items: center;
    justify-content: center;

    height: 33px;
    background: var(--bg-background);
    border-bottom: solid 1px var(--border-gray);
  }

  .panel-tab {
    @include m.flex-row;
    @include m.flex-center-col;

    padding-top: 3px;
    padding-right: 4px;
    margin-bottom: -1px;
    font-size: 14px;
    line-height: 28px;
    color: var(--gray);
    cursor: pointer;
    border-bottom: solid 2px transparent;

    &:hover {
      color: var(--dark-gray);
    }

    &.tab-active {
      border-bottom-color: var(--black);
      color: var(--black);
    }

    .icon {
      width: 16px;
      height: 16px;
      margin-right: 6px;
    }

    .icon-add {
      margin-top: -6px;
    }

    .icon-layer {
      width: 20px;
      height: 20px;
    }
  }

  .panel-tab + .panel-tab {
    margin-left: 20px;
  }

  .panel-tabs-body {
    @include m.flex-col;
    flex: 1 1 0;
  }
</style>
