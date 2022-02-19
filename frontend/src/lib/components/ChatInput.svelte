<script>
    export let value = ""
</script>

<input class="chat-input" type="text" placeholder="메시지 보내기" bind:value={value} on:keypress>

<style lang="scss">
    .chat-input{
        padding: 5px;
        border: solid 1px #c0c0c0;
        border-radius: 3px;
        transition: border 100ms ease-in, background-color 100ms ease-in;
        background: #d0d0d0;
    &::placeholder {
         color: #494949;
     }
    &:focus{
    &::placeholder {
         color: #888888;
     }
    color: #696969;
    outline: none;
    border: solid 1px #772ce8;
    background: transparent;
    }
    }
</style>