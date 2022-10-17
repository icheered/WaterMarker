<script>
  export let text, callback;
</script>

<div>
  <button on:click={callback}>{text}</button>
</div>

<style>
  button {
    border-radius: 50px;
    padding: 10px;
    margin: 5px;
    border: none;
    font-size: 17px;
    cursor: pointer;
    background: #1caeee;
    color: white;
    transition-duration: 0.3s;
    transition-property: background-color;
  }
  button:hover {
    background: #1583b3;
    color: white;
  }
</style>
