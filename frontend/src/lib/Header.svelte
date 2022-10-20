<script>
  export let version, authorLink, author;
</script>

<div class="topbar">
  <div class="author">
    Made by <a  target="_blank" rel="noreferrer" href={authorLink}>{author}</a>
  </div>
  <div class="version">
    {version}
  </div>
  
</div>

<style>
  .topbar {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    text-align: left;
    height: 30px;
    color: black;
    padding: 5px;
  }
  a {
    color: #1caeee;
    text-decoration: none;
  }
</style>
