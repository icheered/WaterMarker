<script>
  export let title, version, authorLink, author;
</script>

<div class="topbar">
  <div class="title">
    {title}
  </div>
  <div class="version">
    {version}
  </div>
  <div class="author">
    <a href={authorLink}>{author}</a>
  </div>
</div>
<div class="line" />

<style>
  .topbar {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    height: 50px;
    padding: 20px;
    background: darkslategrey;
  }
  .line {
    border-bottom: 1px solid white;
    width: 100%;
    margin-bottom: 30px;
  }
  a {
    color: white;
  }
</style>
