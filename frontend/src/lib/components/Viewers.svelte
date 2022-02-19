<script>
    import {format} from 'date-fns'
    import VirtualList from 'svelte-tiny-virtual-list';
    import {onMount} from "svelte";

    let virtualList;
    let scrollToIndex = -1;
    let chatList = [];
    let chatListSize = [];
    let recomputedSizeIndex = 0
    let chatText = ""
    let scrolling = false;

    onMount(() => {

    })

    export function join(data) {
        chatList[chatList.length] = data
        virtualList.recomputeSizes(recomputedSizeIndex);
        recomputedSizeIndex = chatList.length - 1;
    }
    export function leve(data) {
        chatList[chatList.length] = data
        virtualList.recomputeSizes(recomputedSizeIndex);
        recomputedSizeIndex = chatList.length - 1;
    }

    function handleMessage({detail: {startIndex, stopIndex}}) {
        if (!scrolling) {
            scrollToIndex = stopIndex;
        }
    }

    function handleScrollEvent({detail: {event, offset}}) {
        let element = event.target;
        let isBottom = element.scrollHeight - element.scrollTop === element.clientHeight
        if (isBottom) {
            scrolling = false
            scrollToIndex = chatList.length - 1;
            return
        }
        scrolling = true
    }

    function handleKeypress(event) {
        if (event.keyCode === 13) {
            window.go.main.App.SendChatMessage(chatText)
            chatText = ""
        }
    }
</script>

<div class="list">
    <VirtualList
            bind:this={virtualList}
            width="auto"
            height={500}
            itemCount={chatList.length}
            itemSize={30}
            scrollToAlignment={"end"}
            scrollToIndex={(scrolling || scrollToIndex === -1) ? undefined:scrollToIndex}
            getKey={(index)=>{
                console.log("getKey",index)
                return chatList[index].id
            }}
            on:afterScroll={handleScrollEvent}
            on:itemsUpdated={handleMessage}>
        <div slot="item" let:index let:style {style} class="row" class:highlighted={index === scrollToIndex}>
            <div class="message-row"
                 style="position: relative;
                    display: inline-flex;
                    justify-content: flex-start;
                    align-items: center;
                    vertical-align: middle;">
                <span class="time">{format(chatList[index].time, "hh:mm:ss")} </span>
                <div>
                        <span class="badges">
                        {#if chatList[index].badges && chatList[index].badges.length > 0}
                            {#each chatList[index].badges as badge}
                                <img class="badge" src="{badge}" alt=""/>
                            {/each}
                        {/if}
                    </span>
                    <span class="name" style="color: {chatList[index].color};">
                        {chatList[index].display_name !== chatList[index].user_name ? chatList[index].display_name + '(' + chatList[index].user_name + ')' : chatList[index].display_name}
                    </span>
                    {chatList[index].msg}
                </div>
            </div>
        </div>
    </VirtualList>
    {#if scrolling}
        <div class="chat-paused-footer">
            <button>스크롤해서 채팅이 멈췄습니다.</button>
        </div>
    {/if}
</div>
<style>
    .list {
        position: relative;
        min-width: 500px;
    }

    .message-row {
        overflow-wrap: anywhere;
        word-break: break-all;
        height: 100%;
        width: 100%;
    }

    .badges img {
        vertical-align: middle;
        padding: 0;
        margin: 0;
        border: none;
        /*display: inline-block;*/
    }

    .row .time {
        /*display: inline-block;*/
        margin-right: 5px;
        white-space: nowrap;
    }

    .row .badges {
        padding: 0;
        margin: 0;
        border: none;
        margin-right: 5px;
    }

    .row .name {
        display: inline;
        vertical-align: middle;
        margin-right: 5px;
        white-space: nowrap;
    }

    .row .message {
        display: inline-block;
        vertical-align: middle;
    }

    .chat-paused-footer button {
        flex: 1;
        margin: 10px 30px;
        padding: 5px;
        border-radius: 5px;
        border: solid 1px rgba(0, 0, 0, 0.55);
    }

    .chat-paused-footer {
        display: flex;
        position: absolute;
        bottom: 0;
        width: 100%;
    }

    :global(body), :global(html) {
        height: 100%;
        margin: 0;
        background-color: rgb(249, 249, 249);
    }

    .row .name {
        font-weight: 600;
    }

    .row {
        font-size: 13px;
        border-bottom: 1px solid #eee;
        box-sizing: border-box;
        font-weight: 500;
        background: #fff;
    }

    .row img {
        display: inline;
    }

    .row.highlighted {
        background: #efefef;
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


</style>
