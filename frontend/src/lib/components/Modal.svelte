<script>
    export let message = 'default value';
    export let isPromo = false;
    export let showModal = false;
</script>

{#if showModal}
    <div class="backdrop" class:promo={isPromo} on:click|self>
        <div class="modal">
            <slot name="title" />
            <slot />
        </div>
    </div>
{/if}
<style>
    .backdrop {
        width: 100%;
        height: 100%;
        position: fixed;
        top: 0;
        left: 0;
        background: rgba(0, 0, 0, 0.2);
    }
    .modal {
        padding: 10px;
        border-radius: 5px;
        max-width: 400px;
        margin: 10% auto;
        text-align: center;
        background: white;
    }

    .promo .modal {
        background: crimson;
        color: white;
    }

</style>