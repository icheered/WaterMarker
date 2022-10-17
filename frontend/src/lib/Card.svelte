<script>
  export let backgroundcolor;
</script>

<div class="card" style="--backgroundcolor: {backgroundcolor}">
  <slot />
</div>

<style>
  .card {
    padding: 10px;
    margin: 10px;
    border-radius: 20px;
    background-color: #4d5056;
    box-shadow: 0px 0px 20px #4d5056;
  }
</style>
