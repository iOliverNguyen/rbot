<Router url={url}>
  <div class="root">
    <Route path="/" component={Home}/>
    <Route path="/login" component={Login}/>
    <Route path="/editor"><Frame component={Editor}/></Route>
  </div>
</Router>

<script>
  import {Route, Router} from 'svelte-routing';
  import Frame from './com/frame/Frame.svelte';
  import Login from './routes/auth/Login.svelte';
  import Editor from './routes/editor/Editor.svelte';
  import Home from './routes/home/Home.svelte';

  export let url = ''; // svelte-routing


</script>

<style lang="scss">
  .root {
    height: 100%;
    display: flex;
  }
</style>
