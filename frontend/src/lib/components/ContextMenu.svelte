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