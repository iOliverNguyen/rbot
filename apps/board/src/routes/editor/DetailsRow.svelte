<div class="list-item details-row">
  <div class="label">Message</div>
  <div class="content">
    <textarea class="textarea" cols="30" rows="3"></textarea>
  </div>
</div>

<script lang="ts">

</script>

<style lang="scss">
  @use '../../styles/mixins' as m;

  .list-item {
    @include m.flex-col;

    padding: 10px;
    border-bottom: solid 1px var(--washed-gray);
    cursor: pointer;
  }

  .label {
    font-size: 12px;
    font-weight: 600;
    margin-bottom: 6px;
  }

  .textarea {
    width: 100%;
  }
</style>
