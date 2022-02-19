<script>
    import Fa from 'svelte-fa'
    import { faXmark, faPlus, faMinus} from '@fortawesome/free-solid-svg-icons'
</script>

<div class="buttons">
    <button class="close" on:click={()=>{window.runtime.Quit()}}>
        <Fa icon={faXmark} scale={0.85}/>
    </button>
    <button class="minimize">
        <Fa icon={faMinus} scale={0.85}/>
    </button>
    <button class="zoom">
        <Fa icon={faPlus} scale={0.85}/>
    </button>
</div>

<style lang="scss">
  .close {
    background: #fb5e57;
    &:active{
      background: #fb5e57;
    }
  }

  .minimize {
    background: #fbbc2d;
    &:active {
      background: #fbbc2d;
    }
  }

  .zoom {
    background: #28c740;
    &:active {
      background: #28c740;
    }
  }


  .buttons {
    padding-left: 10px;
    padding-right: 10px;
    padding-top: 13px;
    &:hover button{
      visibility: visible;
    }
    button {
      -webkit-appearance: none;
      -moz-appearance: none;
      appearance: none;
      position: relative;
      border-radius: 50%;
      padding: 0;
      display: inline-block;
      width: 13px;
      height: 13px;
      margin: 2px;
      border: none;
    }
    button svg {
      position: absolute;
      left: 0;
      top: 0;
    }
  }

</style>