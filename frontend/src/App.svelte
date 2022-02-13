<script>
    import {onMount} from 'svelte';
    import Chat from "./lib/Chat.svelte";
    let chatBox1;
    let chatBox2;
    let index = 0;
    let b = false;
    onMount(() => {
        window.runtime.EventsOn("chat.stream", (data) => {
            data.time = Date.parse(data.time)
            chatBox1.add(data)
        })
    })

    function greet() {
        window.go.main.App.OpenAuthorization()
    }
</script>

<main>
    <h1>DOVY</h1>
    <div style="display: flex">
        <Chat style="flex:1" bind:this={chatBox1}/>
        <Chat style="flex:1" bind:this={chatBox2} showTime={false}/>
        <div>
            시청자
        </div>
    </div>
    <div>
        저는
        <button>
            스트리머 입니다.
        </button>
        <button>
            매니저 입니다.
        </button>
        <button>
            시청자 입니다.
        </button>
    </div>

    <button class="button" on:click={greet}>인증</button>

</main>

<style>
    :global(.virtual-list-wrapper) {
        margin: 20px;
        background: #fff;
        border-radius: 2px;
        box-shadow: 0 2px 2px 0 rgba(0, 0, 0, .14), 0 3px 1px -2px rgba(0, 0, 0, .2), 0 1px 5px 0 rgba(0, 0, 0, .12);
        background: #fafafa;
        font-family: -apple-system, BlinkMacSystemFont, Helvetica, Arial, sans-serif;
        color: #333;
        -webkit-font-smoothing: antialiased;
    }

    @import url(https://fonts.googleapis.com/earlyaccess/notosanskr.css);

    body, talbe, th, td, div, dl, dt, dd, ul, ol, li, h1, h2, h3, h4, h5, h6,
    pre, form, fieldset, textarea, blockquote, span, * {
        font-family: 'Noto Sans KR', sans-serif;
    }

    .actions label {
        padding: 10px 0;
        font-size: 18px;
        color: #999;
        font-family: -apple-system, BlinkMacSystemFont, Helvetica, Arial, sans-serif;
    }

    .input:focus {
        border-bottom: 2px solid #008cff;
        margin-bottom: 9px;
    }


    :global(input::-moz-focus-inner), :global(input::-moz-focus-outer) {
        border: 0;
    }
</style>
