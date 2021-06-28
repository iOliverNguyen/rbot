<div class="auth-root">
  <div class="auth-box">
    <div class="auth-header">
      <div class="logo"></div>
      <div class="title">ReviewBot Login</div>
    </div>
    <div class="auth-text">
      <div class="header">Welcome back to ReviewBot Board!</div>
      <div class="message"></div>
    </div>
    <div class="auth-form">
      <div class="auth-provider">
        <div class="auth-button btn-button">
          <img class="icon" src={logoG}>
          <span>Sign in with Google</span>
        </div>
      </div>
      <div class="hline">
        <div class="or">or</div>
      </div>
      <form class="auth-form-input">
        <div class="row">
          <input type="text" class="text" placeholder="Email">
        </div>
        <div class="row">
          <input type="password" class="password" placeholder="Password">
        </div>
        <div class="row row-h">
          <label class="h-item clickable">
            <input type="checkbox" class="checkbox clickable">
            <span class="text">Remember me</span>
          </label>
          <div class="h-item forgot-password">
            <a class="link" href="/forgot-password">Forgot password?</a>
          </div>
        </div>
      </form>
      <div class="auth-buttons">
        <div class="btn-button">Sign Up</div>
        <div class="btn-button btn-primary">Sign In</div>
      </div>
      <div class="footer-text">
        By using ReviewBot, you are agreeing to our
        <a href="/doc/terms">terms and conditions</a>.
      </div>
    </div>
  </div>
</div>

<script lang="ts">
  import logoG from '../../assets/logo-g.svg';
</script>

<style lang="scss">
  @use "../../styles/mixins" as m;

  .auth-root {
    @include m.flex-row;
    @include m.flex-center;

    background: var(--bg-background);
    flex: 1 0 auto;
  }

  .auth-box {
    @include m.flex-col;

    width: 400px;
    background: var(--bg-dialog-box);
    padding: 48px 30px 48px;
    border-radius: 6px;
    border: 1px solid var(--washed-gray);
    box-shadow: 0 3px 6px rgb(0 0 0 / 10%);
  }

  .auth-header {
    @include m.flex-row;
    @include m.flex-center;

    .logo {
      width: 40px;
      height: 40px;
      background: url(/_/assets/icon-purple.png) center no-repeat;
      background-size: contain;
    }

    .title {
      margin-left: 20px;
      font-size: 20px;
      font-weight: 600;
    }
  }

  .auth-text {
    @include m.flex-row;
    @include m.flex-center-row;

    margin: 30px 0 40px 0;
  }

  .auth-form {
    @include m.flex-col;
  }

  .auth-provider {
    @include m.flex-row;
    @include m.flex-center-row;
  }

  .auth-button {
    @include m.flex-row;
    @include m.flex-center;

    flex: 1 0 auto;

    .icon {
      margin-right: 10px;
      width: 20px;
    }
  }

  .hline {
    @include m.flex-row;
    @include m.flex-center;

    height: 1px;
    width: 100%;
    background: var(--washed-gray);
    margin: 20px 0;

    .or {
      background: var(--bg-dialog-box);
      padding: 0 20px;
      color: var(--light-gray)
    }
  }

  .auth-form-input {
    @include m.flex-col;

    .row {
      @include m.flex-row;

      flex: 1 0 auto;
      width: 100%;
      margin-top: 10px;
    }

    input {
      flex: 1 0 auto;
      display: block;
      font-size: 14px;
      padding: 9px 10px;
      border-radius: 4px;
      border: solid 1px var(--washed-gray);
    }

    .checkbox {
      width: 14px;
      height: 14px;
      margin-right: 6px;
    }

    .row-h {
      @include m.flex-row;
      justify-content: space-between;
      font-size: 14px;
    }

    .h-item {
      display: flex;
      align-items: center;
    }
  }

  .auth-buttons {
    @include m.flex-row;
    justify-content: center;

    margin-top: 30px;

    .btn-button {
      flex: 1 0 0;
    }

    .btn-button + .btn-button {
      margin-left: 20px;
    }
  }

  .footer-text {
    margin-top: 16px;
    font-size: 12px;
    color: var(--gray);
    text-align: center;
  }
</style>
