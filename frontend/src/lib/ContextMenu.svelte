<script>
    let ctxMenu
    export let menus = []

    export function Close() {
        ctxMenu.style.display = 'none';
        ctxMenu.style.top = null;
        ctxMenu.style.left = null;
    }

    export function Open(position, event) {
        console.log(position, event)
        event.preventDefault();
        ctxMenu.style.display = 'block';
        let ch = document.documentElement.clientHeight

        if (position === 'left') {
            ctxMenu.style.left = event.pageX - ctxMenu.offsetWidth + 'px';
            console.log(event.pageX, event.pageY, ctxMenu.offsetWidth,ch)
        } else {
            ctxMenu.style.left = event.pageX + 'px';
        }
        if (ch < event.pageY + ctxMenu.offsetHeight) {
            ctxMenu.style.top = event.pageY - ctxMenu.offsetHeight + 'px';
        } else {
            ctxMenu.style.top = event.pageY + 'px';
        }
    }
</script>

<div bind:this={ctxMenu} class="custom-context-menu" style="display: none;">
    <slot/>
</div>

<style>
    .custom-context-menu ul li:last-child {
        border-left: 0;
        border-right: 0;
        border-top: 0;
        border-bottom: 0;
    }

    ul {
        list-style-type: none;
        padding: 0;
        margin: 0;
    }

    .custom-context-menu ul li {
        cursor: pointer;
        min-width: 130px;
        font-size: 12px;
        padding: 5px;
        text-decoration: none;
        border-bottom: 1px #C5C5C5FF;
        border-left: 0;
        border-right: 0;
        border-top: 0;
        border-style: solid;
    }

    .custom-context-menu ul li:hover {
        background-color: #d7d7d7;
    }

    .custom-context-menu {
        max-width: 240px;
        font-family: 'Noto Sans KR', sans-serif;
        border: solid 1px #C5C5C5FF;
        border-radius: 3px;
        position: absolute;
        box-sizing: border-box;
        background-color: #ffffff;
        box-shadow: 0 0 1px 2px lightgrey;
        -webkit-box-shadow: 0px 0px 5px 3px rgba(0, 0, 0, 0.1);
    }
</style>