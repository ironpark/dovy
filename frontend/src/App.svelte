<script>
    import {onMount} from 'svelte';
    import Chat from "./lib/Chat.svelte";
    import Fa from 'svelte-fa'
    import {faEllipsisVertical, faGear, faGears, faXmark, faPlus, faMinus} from '@fortawesome/free-solid-svg-icons'
    import Tabs from "./lib/Tabs.svelte";
    import Modal from "./lib/Modal.svelte";

    let showModal = false;
    let showLoginModal = false;

    let chatBox1;
    let chatBox2;
    let index = 0;
    let b = false;
    let selectedData = [];
    let selectedIndex = 0;
    let chatMessage = ""
    let channelName = ""
    let channels = {};
    let selectedChannel = {
        filter: {},
        name: '',
        itemList: [],
        users: {},
    }
    let tabs = [];
    let tabData = {};
    onMount(async () => {
        showLoginModal = !await window.go.main.App.IsAuthorized()
        window.runtime.EventsOn("stream.chat", async (data) => {
            data.time = Date.parse(data.time)
            let {channel, created} = channelInitIfNotExist(data.channel)
            if (created) {
                let userlist = await window.go.main.App.UserList(channel.name);
                for (let i = 0; i < userlist.length; i++) {
                    channel.users[userlist[i]] = true;
                }
                selectedIndex = tabs.length - 1;
                selectedChannel = channel;
            }

            channel.itemList[channel.itemList.length] = data;

            if (channel.itemList.length > 1000) {
                channel.itemList = channel.itemList.slice(channel.itemList.length - 1 - 110, channel.itemList.length - 1)
            }

            if (selectedChannel.name === data.channel) {
                selectedChannel = channel;
            }
        })
        window.runtime.EventsOn("stream.user-event", (data) => {
            let {channel} = channelInitIfNotExist(data.channel)
            if (data.event === "JOIN") {
                if (channel.users[data.user] === undefined) {
                    channel.users[data.user] = true
                }
            } else {
                if (channel.users[data.user] === undefined) {
                    return
                }
                delete channel.users[data.user];
            }
            console.log(data.channel, data.event, data.user)
        })
    })

    function channelExist(channelName) {
        return channels[channelName] !== undefined
    }

    function channelInitIfNotExist(channelName) {
        if (channels[channelName] !== undefined) {
            return {channel: channels[channelName], created: false}
        }
        let channel = {
            filter: {},
            name: channelName,
            itemList: [],
            users: {}
        }
        channels[channelName] = channel;

        let tabItem = {
            name: channelName,
            views: [{
                name: channelName,
                filter: {},
                channel: channel,
            }]
        }

        tabs[tabs.length] = channelName;
        tabData[channelName] = tabItem;
        return {channel, created: true}
    }

</script>
<Modal bind:showModal={showModal}>
    <div slot="title">
        <h3>채팅 채널 추가</h3>
    </div>
    <input type="text" placeholder="스트리머 아이디" bind:value={channelName}/>
    <button on:click={async ()=>{
            if(channelName === "") return;
            if(channelExist(channelName))return;
            channelName = channelName.trim()
            let {channel,created} = channelInitIfNotExist(channelName);
            await window.go.main.App.Connect(channelName);
            let userlist = await window.go.main.App.UserList(channel.name);
            if (userlist){
                for(let i = 0; i < userlist.length; i++) {
                    channel.users[userlist[i]] = true;
                }
            }
            selectedIndex = tabs.length - 1;
            selectedChannel = channel
            showModal = false;
        }}> 추가
    </button>
    <button on:click={()=>{
            showModal = false;
            channelName = "";
    }}>취소
    </button>
</Modal>

<Modal bind:showModal={showLoginModal}>
    <div slot="title">
        <h3>로그인</h3>
    </div>
    <button on:click={()=>{
            window.go.main.App.OpenAuthorization()
            showLoginModal = false
        }}>트위치 로그인
    </button>
</Modal>

<main>
    <Tabs bind:selectedIndex={selectedIndex} bind:tabs={tabs} on:addclick={()=>{showModal = true}}
          on:selectTab={(e)=>{
              selectedChannel = channels[tabs[e.detail.index]]
              selectedIndex = e.detail.index
              chatBox1.scrollToBottom();
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

    <div style="display: flex;height:calc(100vh - 70px);flex-flow: row">
        <Chat style="flex:1;height: 100%;" bind:this={chatBox1} bind:itemList={ selectedChannel.itemList }/>
        <div class="user-list">
            {#each Object.keys(selectedChannel.users) as user}
                <div>{user}</div>
            {/each}
        </div>
    </div>
    <input style="" type="text" bind:value={chatMessage} on:keypress={(e)=>{
            if (e.charCode === 13 && (chatMessage && chatMessage.trim() !== ""))  {
                window.go.main.App.SendChatMessage(selectedChannel.name,chatMessage)
                chatMessage = ""
            }
        }}>

</main>

<style>
    main {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
    }

    :global(html,body) {
        margin: 0;
        height: 100%;
    }

    :global(#app) {
        height: 100%;
        display: flex;
    }

    .contents {
        display: flex;
        flex-direction: column;
        width: 100%;
        overflow-y: auto;
    }

    @import url(https://fonts.googleapis.com/earlyaccess/notosanskr.css);
    body, talbe, th, td, div, dl, dt, dd, ul, ol, li, h1, h2, h3, h4, h5, h6,
    pre, form, fieldset, textarea, blockquote, span, * {
        font-family: 'Noto Sans KR', sans-serif;
    }

    :global(input::-moz-focus-inner), :global(input::-moz-focus-outer) {
        border: 0;
    }

    .user-list {
        width: 130px;
        overflow-y: scroll;
        font-size: 13px;
        padding-left: 5px;
        border-left: solid 1px #999999;
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
        border: none;
    }

    .buttons button svg {
        position: absolute;
        left: 0;
        top: 0;
    }

    .close {
        background: #fb5e57;
    }

    .close:active {
        background: #fb5e57;
    }

    .minimize {
        background: #fbbc2d;
    }

    .minimize:active {
        background: #fbbc2d;
    }

    .zoom {
        background: #28c740;
    }

    .zoom:active {
        background: #28c740;
    }


</style>
