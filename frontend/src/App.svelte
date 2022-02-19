<script>
    import {onMount} from 'svelte';
    import Chat from "./lib/components/Chat.svelte";
    import Tabs from "./lib/components/Tabs.svelte";
    import Modal from "./lib/components/Modal.svelte";
    import ContextMenu from "./lib/components/ContextMenu.svelte";
    import ChatInput from "./lib/components/ChatInput.svelte";
    import WindowButtonGroup from "./lib/components/WindowButtonGroup.svelte";

    let showModal = false;
    let showLoginModal = false;
    let ctxMenu;
    let chatCtxMenu;
    let chatBox1;

    let selectedIndex = 0;
    let selectedMessage = "";

    let chatMessage = "";
    let channelName = "";
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

<main on:click={(e)=>{ctxMenu.Close();chatCtxMenu.Close()}}>
    <Tabs bind:selectedIndex={selectedIndex} bind:tabs={tabs} on:addclick={()=>{showModal = true}}
          on:selectTab={(e)=>{
              selectedChannel = channels[tabs[e.detail.index]]
              selectedIndex = e.detail.index
              chatBox1.scrollToBottom();
          }}>
        <div slot="front">
            <WindowButtonGroup/>
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
                <div on:contextmenu={(e)=>{selectedUser = user;ctxMenu.Open('left',e)}}>{user}</div>
            {/each}
        </div>
    </div>
    <ChatInput  bind:value={chatMessage} on:keypress={(e)=>{
            if (e.charCode === 13 && (chatMessage && chatMessage.trim() !== ""))  {
                window.go.main.App.SendChatMessage(selectedChannel.name,chatMessage)
                chatMessage = ""
            }
        }}/>
</main>

<ContextMenu bind:this={ctxMenu}>
    <div style="font-size: 13px;padding: 5px;border-bottom: solid 1px #a2a2a2"> @{selectedUser} </div>
</ContextMenu>

<ContextMenu bind:this={chatCtxMenu}>
    <div style="font-weight: bold;font-size: 13px;min-width: 100px">@{selectedUser}</div>
    <div style="font-size: 12px;padding: 3px;">
        {@html selectedMessage}
    </div>
    <div style="display: flex;flex-flow: row;margin: 0;padding: 0;font-size: 13px;" class="menu-item">
        <div class="menu-item" style="width: 30px;text-align: center">답변</div>
        <div class="menu-item" style="width: 30px;text-align: center">맨션</div>
        <div class="menu-item" style="width: 80px;text-align: center">내용 복사</div>
        <div class="menu-item" style="width: 90px;text-align: center">아이디 복사</div>
    </div>
    <div style="display: flex;flex-flow: column;margin: 0;padding: 0;font-size: 13px;">
        <div class="menu-item">임시차단</div>
        <div class="menu-item">강퇴</div>
    </div>
</ContextMenu>

<Modal bind:showModal={showModal}>
    <div slot="title">
        <h3>채팅 채널 추가</h3>
    </div>
    <input type="text" placeholder="스트리머 아이디" bind:value={channelName}/>
    <button on:click={async ()=>{
            if(channelName === "") return;
            if(channelExist(channelName))return;
            channelName = channelName.trim()
            let {channel} = channelInitIfNotExist(channelName);
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

<style lang="scss">
  :global {
    @import "global";
  }

  .contents {
    overscroll-behavior: none;
    display: flex;
    flex-direction: column;
    width: 100%;
    overflow-y: auto;
  }

  .user-list div:hover {
    background-color: #e3e3e3;
  }

  .menu-item {
    .menu-item {
      border-right: solid 1px silver;
      border-top: 0;
    }
    &:last-child {
      border-right: 0;
    }
    border-top: solid 1px silver;
    padding: 3px;
  }

  .user-list {
    width: 130px;
    overflow-y: scroll;
    font-size: 13px;
    border-left: solid 1px #999999;
    overflow-x: hidden;
    div {
      -webkit-user-select: none;
      padding-top: 2px;
      padding-bottom: 2px;
      padding-left: 5px;
    }
  }

  span {
    line-height: 9px;
    vertical-align: 50%;
  }
</style>
