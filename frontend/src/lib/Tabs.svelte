<script>
    import Fa from 'svelte-fa/src/fa.svelte'
    import {faEllipsisVertical, faGear, faGears, faXmark, faPlus} from '@fortawesome/free-solid-svg-icons'
    import {createEventDispatcher} from 'svelte';

    const dispatch = createEventDispatcher();
    export let tabs = [];
    export let selectedIndex = 0;

    function subTabAdd(index) {
        dispatch('sub', {
            index: index,
        });
    }

    function select(index) {
        dispatch('select', {
            index: index,
        });
    }

    function addClick() {
        dispatch('addclick', {});
    }
</script>

<div class="tabs">
    <div class="tab-container">
        <div data-wails-drag>
            <slot name="front"/>
        </div>
        {#each tabs as tab, i}
            <div class="tab {i===selectedIndex?'active':''}"
                 on:mousedown|self={(e)=>{
                     select(i);
                     selectedIndex = i;

                     e = e || window.event;
                     let start = 0, diff = 0;
                     start = e.pageX ? e.pageX:(e.clientX?e.clientX:0)

                     e.target.style.position = 'relative';
                     document.body.onmousemove = function(e) {
                            e = e || window.event;
                            let end = e.pageX ? e.pageX:(e.clientX?e.clientX:0);
                            diff = end-start;
                            e.target.style.left = diff+"px";
                     };
                     document.body.onmouseup = function() {
                         // do something with the action here
                         // elem has been moved by diff pixels in the X axis
                         e.target.style.left = "0px";
                         e.target.style.position = 'static';
                         document.body.onmousemove = document.body.onmouseup = null;
                     };
                 }}>
                {tab}
                <button style="padding-left: 5px;padding-right: 5px;color: #999999" on:click={()=>subTabAdd(i)}>
                    <Fa icon={faEllipsisVertical} pull="left" flip="horizontal" scale={1.2}/>
                </button>
                <button class="close" on:click={()=>subTabAdd(i)}>
                    <Fa icon={faXmark} pull="left" flip="horizontal" scale={1.2}/>
                </button>
            </div>
            {#if i !== selectedIndex && i !== selectedIndex - 1}
                <div class="divider"></div>
            {:else}
                <div class="divider tp"></div>
            {/if}
        {/each}

        <div class="add-tab-button">
            <button on:click={addClick}>
                <Fa icon={faPlus} pull="left" flip="horizontal" scale={1.2}/>
            </button>
        </div>
        <div class="tab-spacer" data-wails-drag>
        </div>
        <slot name="back"/>
    </div>
</div>
<style>

    .tabs {
        width: 100%;
        display: flex;
        flex-direction: row;
        background: #cecece;
        -webkit-touch-callout: none; /* iOS Safari */
        -webkit-user-select: none; /* Safari */
        -khtml-user-select: none; /* Konqueror HTML */
        -moz-user-select: none; /* Old versions of Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
        user-select: none;
        /* Non-prefixed version, currently
                                         supported by Chrome, Edge, Opera and Firefox */
    }

    .tab-container {
        width: 100%;
        height: 41px;
        display: flex;
        flex-direction: row;
    }

    .tabs .tab {
        font-size: 15px;
        padding: 12px 10px;
        background: transparent;
    }

    .add-tab-button {
    }

    .tab button {
        -webkit-appearance: none;
        -moz-appearance: none;
        appearance: none;
        border: none;
        display: inline-block;
        margin: 0;
        padding: 3px;
        background: transparent;
    }

    .tab button.close {
        background: transparent;
    }

    .add-tab-button button {
        -webkit-appearance: none;
        -moz-appearance: none;
        appearance: none;
        border: none;
        min-width: 27px;
        min-height: 27px;
        width: 100%;
        height: 100%;
        background: transparent;
        padding-left: 10px;
        padding-right: 10px;
    }

    .tab-spacer {
        flex: 1;
    }

    .tab.active {
        border-bottom: none;
        background: white;
        margin-top: 5px;
        padding-top: 7px;
        padding-bottom: 12px;
        border-top-left-radius: 5px;
        border-top-right-radius: 5px;
        -webkit-box-shadow: 0px -14px 13px 3px rgba(0, 0, 0, 0.2);
        box-shadow: 0px -14px 13px 3px rgba(0, 0, 0, 0.2);
        border-right: none;
    }

    .divider {
        width: 1px;
        margin: 9px 0;
        background: #919191;
    }

    .divider.tp {
        width: 1px;
        margin: 9px 0;
        background: transparent;
    }
</style>
