<script>
    import {onMount} from 'svelte';
    import Chat from "./lib/Chat.svelte";
    import Fa from 'svelte-fa'
    import {faEllipsisVertical, faGear, faGears, faXmark, faPlus, faMinus} from '@fortawesome/free-solid-svg-icons'
    import Tabs from "./lib/Tabs.svelte";
    import Modal from "./lib/Modal.svelte";
    import ContextMenu from "./lib/ContextMenu.svelte";

    let showModal = false;
    let showLoginModal = false;
    let ctxMenu;
    let chatCtxMenu;
    let chatBox1;
    let chatBox2;
    let index = 0;
    let b = false;
    let selectedData = [];
    let selectedIndex = 0;
    let selectedMessage = ""
    let chatMessage = ""
    let channelName = ""
    let channels = {};
    let selectedUser = "";
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

<ContextMenu bind:this={ctxMenu} menus={["맨션하기","아이디 복사하기","임시차단","강퇴"]}>
    <div style="font-size: 13px;padding: 5px;border-bottom: solid 1px #a2a2a2"> @{selectedUser} </div>
</ContextMenu>

<ContextMenu bind:this={chatCtxMenu} menus={["답변하기","맨션하기","내용 복사","아이디 복사","임시차단","강퇴"]}>
    <div style="font-weight: bold;font-size: 13px;min-width: 100px">@{selectedUser}</div>
    <div style="font-size: 12px;padding: 3px;">
        {@html selectedMessage}
    </div>
    <div style="display: flex;flex-flow: row;margin: 0;padding: 0;font-size: 13px;" class="menu-item">
        <div class="menu-item" style="width: 30px;text-align: center">답변</div>
        <div class="menu-item" style="width: 30px;text-align: center">맨션</div>
        <div class="menu-item" style="width: 80px;text-align: center">내용 복사</div>
        <div class="menu-item" style="width: 90px;;text-align: center">아이디 복사</div>
    </div>
    <div style="display: flex;flex-flow: column;margin: 0;padding: 0;font-size: 13px;">
        <div class="menu-item">임시차단</div>
        <div class="menu-item">강퇴</div>
    </div>

</ContextMenu>

<main on:click={(e)=>{
     ctxMenu.Close()
     chatCtxMenu.Close()
}}>

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
        <Chat style="flex:1;height: 100%;" bind:this={chatBox1} bind:itemList={ selectedChannel.itemList }
              callCtxMenu={(e)=>{
                  selectedUser = e.user;
                  selectedMessage = e.message;
                  chatCtxMenu.Open('right',e);
              }}/>
        <div class="user-list">
            {#each Object.keys(selectedChannel.users) as user (user)}
                <div on:contextmenu={(e)=>{
                    selectedUser = user
                    ctxMenu.Open('left',e)
                }}>{user}</div>
            {/each}
        </div>
    </div>
    <input class="chat-input" type="text" placeholder="메시지 보내기" bind:value={chatMessage} on:keypress={(e)=>{
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
        overscroll-behavior: none;
    }

    :global(::-webkit-scrollbar) {
        width: 5px;
    }

    :global(::-webkit-scrollbar-thumb) {
        background-color: grey;
        border: solid 1px transparent;
        border-radius: 3px;
        margin: 1px;
    }

    :global(::-webkit-scrollbar-track) {
        width: 7px;
        background-color: transparent;
    }

    :global(html,body) {
        overscroll-behavior: none;
        margin: 0;
        height: 100%;
        touch-action: none;
        background-attachment: fixed;
    }

    :global(#app) {
        overscroll-behavior: none;
        height: 100%;
        display: flex;
    }


    .contents {
        overscroll-behavior: none;
        display: flex;
        flex-direction: column;
        width: 100%;
        overflow-y: auto;
    }

    .chat-input {
        padding: 5px;
        border: solid 1px #c0c0c0;
        border-radius: 3px;
        transition: border 100ms ease-in, background-color 100ms ease-in;
        background: #d0d0d0;
    }

    .chat-input::placeholder {
        color: #494949;
    }

    .chat-input:focus::placeholder {
        color: #696969;
    }
    .menu-item .menu-item {
        border-right: solid 1px silver;
        border-top:0;
    }

    .menu-item:last-child{
        border-right: 0;
    }

    .menu-item{
        border-top: solid 1px silver;
        padding: 3px;
    }

    .chat-input:focus {
        outline: none;
        border: solid 1px #772ce8;
        background: transparent;
    }
    @import url(https://fonts.googleapis.com/earlyaccess/notosanskr.css);
    :global(body, talbe, th, td, div, dl, dt, dd, ul, ol, li, h1, h2, h3, h4, h5, h6,
    pre, form, fieldset, textarea, blockquote, span, *) {
        font-family: 'Noto Sans KR', sans-serif;
    }

    :global(input::-moz-focus-inner), :global(input::-moz-focus-outer) {
        border: 0;
    }

    .user-list {
        width: 130px;
        overflow-y: scroll;
        font-size: 13px;
        border-left: solid 1px #999999;
        overflow-x: hidden;
    }

    .user-list div {
        -webkit-user-select: none;
        padding-top: 2px;
        padding-bottom: 2px;
        padding-left: 5px;
    }

    .user-list div:hover {
        background-color: #e3e3e3;
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


    .custom-context-menu ul {
        list-style: none;
        padding: 0;
        background-color: transparent;
    }

    .custom-context-menu li {
        padding: 3px 5px;
        cursor: pointer;
    }

    .custom-context-menu li:hover {
        background-color: #f0f0f0;
    }
</style>
