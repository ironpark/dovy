<script>
    import {format} from 'date-fns'
    import VirtualScroll from "svelte-virtual-scroll-list"
    import {afterUpdate, onMount} from "svelte";
    import Fa from "svelte-fa";
    import {faReply} from "@fortawesome/free-solid-svg-icons";

    export let style = "width:300px;";
    export let showBadges = true;
    export let showEmotes = true;
    export let showTime = true;
    export let itemList = [];
    export let callCtxMenu;
    let listElement;
    let virtualList;
    let scrolling = false;
    onMount(() => {

    })

    export function add(data) {
        if (itemList.length > 200) {
            itemList = [...itemList.slice(itemList.length - 1 - 50, itemList.length - 1), data]
        } else {
            itemList[itemList.length] = data
        }

        if (!scrolling) virtualList.scrollToBottom();
    }

    export function scrollToBottom() {
        listElement.scrollTo(0, listElement.scrollHeight);
    }

    afterUpdate(() => {
        listElement.scrollTo(0, listElement.scrollHeight);
    });

    function handleScrollEvent({detail: {event, offset}}) {
        let element = event.target;
        let isBottom = Math.abs(element.scrollHeight - element.scrollTop - element.clientHeight) < 100
        if (isBottom) {
            scrolling = false
            virtualList.scrollToBottom();
            return
        }
        scrolling = true
    }

</script>

<div class="list" style={style} bind:this={listElement} on:copy={(e)=>{
  let selection = document.getSelection(),
  range = selection.getRangeAt(0),
  contents = range.cloneContents(),
  copiedText = '';

  for (let node of contents.childNodes.values()) {
    if (node.nodeType === 3) {
      // text node
      copiedText += node.textContent;
    } else if (node.nodeType === 1 && node.nodeName === 'IMG') {
      copiedText += node.alt;
    }
  }
  e.clipboardData.setData('text/plain', copiedText);
  e.preventDefault();
}}>
    {#each itemList as data (data.id)}
        <div class="message-row"
             style="display: flex;
                    justify-content: flex-start;
                    align-items: start;
                    vertical-align: middle;" on:contextmenu={(e)=>{
                        e.user=data.user_name;
                        e.message = data.msg_with_emotes;
                        callCtxMenu(e)}}>
            {#if showTime}
                <div class="time">{format(data.time, "hh:mm:ss")} </div>
            {/if}
            <div style="flex: 1;">
                {#if showBadges && data.badges && data.badges.length > 0}
                    {#each data.badges as badge}
                        <img class="badge" src="{badge}" alt=""/>
                    {/each}
                {/if}
                <span class="name" style="color: {data.color};">
                        {data.is_user_name_only ? data.display_name : data.display_name + '(' + data.user_name + ')'}
                </span>
                {@html data.msg_with_emotes}
            </div>
            <div class="actions">
                <button>@</button>
                <button>
                    <Fa icon={faReply} scale={0.85}/>
                </button>
            </div>
        </div>
    {/each}
    <!--    <div bind:this={listElement}></div>-->
    <!--    </VirtualScroll>-->
    {#if scrolling}
        <div class="chat-paused-footer">
            <button>스크롤해서 채팅이 멈췄습니다.</button>
        </div>
    {/if}
</div>
<style>
    .list {
        overflow-y: auto;
    }

    :global(.virtual-list-wrapper) {
        margin: 20px;
        border-radius: 2px;
        box-shadow: 0 2px 2px 0 rgba(0, 0, 0, .14), 0 3px 1px -2px rgba(0, 0, 0, .2), 0 1px 5px 0 rgba(0, 0, 0, .12);
        background: #fafafa;
        font-family: -apple-system, BlinkMacSystemFont, Helvetica, Arial, sans-serif;
        color: #333;
        -webkit-font-smoothing: antialiased;
    }

    :global(img.emote) {
        display: inline;
        vertical-align: middle;
        height: 17px;
    }

    :global(.message-row .badge) {
        height: 17px;
    }

    .message-row {
        position:relative;
        min-height: 18px;
        font-size: 13px;
        overflow-wrap: anywhere;
        word-break: break-all;
        padding: 5px 20px 5px 5px;
    }
    .actions{
        position:absolute;
        right: 0;
        top: 3px;
        display: none;
    }
    .actions button{
        text-align: center;
        padding: 4px;
        margin: 0;
    }
    .message-row:hover .actions{
        display: block;
    }
    .message-row .name:hover {
        background: #d2d2d2;
    }

    .message-row:hover {
        background: #e3e3e3;
    }

    .name {
        padding: 2px;
        cursor: pointer;
    }

    .message-row .badge {
        object-fit: none;
        vertical-align: middle;
        padding: 0;
        margin: 0 2px 0 0;
        border: none;
    }

    .message-row .time {
        margin-right: 5px;
        white-space: nowrap;
    }

    .message-row .name {
        display: inline;
        vertical-align: middle;
        margin-right: 5px;
        white-space: nowrap;
        font-weight: 600;
    }

    .chat-paused-footer button {
        flex: 1;
        margin: 10px 30px;
        padding: 5px;
        border-radius: 5px;
        border: solid 1px rgba(0, 0, 0, 0.55);
    }

    .chat-paused-footer {
        position: absolute;
        bottom: 0;
        width: 100%;
    }

</style>
