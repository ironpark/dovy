<script>
    import {onMount} from 'svelte';
    import Chat from "./lib/Chat.svelte";
    import Fa from 'svelte-fa'
    import {faEllipsisVertical, faGear, faGears, faXmark, faPlus, faMinus} from '@fortawesome/free-solid-svg-icons'
    import Tabs from "./lib/Tabs.svelte";
    import Modal from "./lib/Modal.svelte";

    let showModal = false;
    let chatBox1;
    let chatBox2;
    let index = 0;
    let b = false;
    let tabs = ["#asmongold", "#vo_ine"];
    let selectedIndex = 0;

    let channelName = ""
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
    <Tabs bind:selectedIndex={selectedIndex} bind:tabs={tabs} on:addclick={()=>{
        showModal = true
        console.log("?")
    }}>
        <div slot="front">
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
        </div>
    </Tabs>
    <div class="contents">
        <div style="display: flex">
            <Chat style="flex:1" bind:this={chatBox1}/>
            <Chat style="flex:1" bind:this={chatBox2} showTime={false}/>
            <div>
                시청자
            </div>
        </div>
    </div>
    <button class="button" on:click={greet}>인증</button>
</main>

<Modal bind:showModal={showModal}>
    <div slot="title">
        <h3>채팅 채널 추가</h3>
    </div>
    <input type="text" placeholder="스트리머 아이디" bind:value={channelName}/>
    <button on:click={()=>{
            console.log("!!",channelName)
            tabs = [...tabs,"#"+channelName]
            channelName = ""
            showModal = false
            selectedIndex = tabs.length - 1;
        }}>추가
    </button>
    <button>취소</button>
</Modal>

<style>

    :global(body) {
        margin: 0;
    }

    @import url(https://fonts.googleapis.com/earlyaccess/notosanskr.css);
    body, talbe, th, td, div, dl, dt, dd, ul, ol, li, h1, h2, h3, h4, h5, h6,
    pre, form, fieldset, textarea, blockquote, span, * {
        font-family: 'Noto Sans KR', sans-serif;
    }

    :global(input::-moz-focus-inner), :global(input::-moz-focus-outer) {
        border: 0;
    }
    /* window BEGIN */

    span {
        line-height: 9px;
        vertical-align: 50%;
    }

    .buttons {
        padding-left: 10px;
        padding-right: 10px;
        padding-top: 13px;
    }

    .buttons:hover button {
        visibility: visible;
    }

    .buttons button {
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
    }
    .buttons button svg{
        position: absolute;
        left: 0;
        top: 0;
    }
    .close {
        background: #ff5c5c;
        border: 1px solid #e33e41;
    }
    .close:active {
        background: #c14645;
        border: 1px solid #b03537;
    }
    .close:active .closebutton {
        color: #4e0002;
    }

    .minimize {
        background: #ffbd4c;
        border: 1px solid #e09e3e;
    }

    .minimize:active {
        background: #c08e38;
        border: 1px solid #af7c33;
    }

    .minimize:active .minimizebutton {
        color: #5a2607;
    }

    .zoom {
        background: #00ca56;
        border: 1px solid #14ae46;
    }

    .zoom:active {
        background: #029740;
        border: 1px solid #128435;
    }

    .zoom:active .zoombutton {
        color: #003107;
    }

</style>
