<script>
  import { ProcessImages, FetchPreview } from "../wailsjs/go/main/App.js";
  import Header from "./lib/Header.svelte";
  import Pathconfig from "./lib/Pathconfig.svelte";
  import Settings from "./lib/Settings.svelte";
  import Card from "./lib/Card.svelte";
  import Button from "./lib/Button.svelte";
  import Loader from "./lib/Loader.svelte";

  // CONSTANTS
  const title = "WaterMarker";
  const version = "1.2";
  const author = "ICheered";
  const authorLink = "https://icheered.nl";

  // PARAMETERS
  // let sourceFolderPath = "/home/tjbakker/Documents/dev/vscode/go/WaterMarker/testfiles/source";
  // let targetFolderPath = "/home/tjbakker/Documents/dev/vscode/go/WaterMarker/testfiles/watermarked";
  // let watermarkPath = "/home/tjbakker/Documents/dev/vscode/go/WaterMarker/testfiles/watermark.png";

  let sourceFolderPath = "";
  let targetFolderPath = "";
  let watermarkPath = "";

  let watermarkOpacity = 80;
  let watermarkPosition = "bottom-right";
  let watermarkScale = 20;

  // INFORMATION
  let changedSettings = true;

  let numberOfSourceFiles = 0;
  let numberOfTargetFiles = 0;

  let watermarkpreviewImage;
  let watermarkedpreviewImage;

  let showLoader = false;

  // FUNCTIONS
  let mainbuttontext = "";
  $: {
    if (sourceFolderPath == "" || targetFolderPath == "" || watermarkPath == "") {
      mainbuttontext = "Select source, target and watermark";
    } else if (numberOfSourceFiles == 0) {
      mainbuttontext = "No images found in source folder!";
    } else if (changedSettings) {
      mainbuttontext = "Generate preview";
    } else {
      mainbuttontext = "Start processing!";
    }
  }

  function generatePreview() {
    // fetch("/home/tjbakker/Documents/dev/vscode/go/wails/testfiles/watermarked/DSC_0134.jpg")
    //   .then((response) => response.blob())
    //   .then((blob) => {
    //     const reader = new FileReader();
    //     reader.addEventListener("load", function () {
    //       watermarkedpreviewImage.setAttribute("src", reader.result);
    //     });
    //     reader.readAsDataURL(blob);
    //   });
    showLoader = true;
    console.log("Generating preview");
    FetchPreview(
      watermarkPath,
      sourceFolderPath,
      targetFolderPath,
      watermarkPosition,
      watermarkOpacity / 100,
      watermarkScale / 100
    ).then((result) => {
      if (result.status && result.status == "error") {
        alert(result.message);
      } else if (result.status && result.status == "success") {
        fetch(result.message)
          .then((response) => response.blob())
          .then((blob) => {
            const reader = new FileReader();
            reader.addEventListener("load", function () {
              watermarkedpreviewImage.setAttribute("src", reader.result);
            });
            reader.readAsDataURL(blob);
            showLoader = false;
            changedSettings = false;
          });
      }
    });
  }

  function processFiles() {
    showLoader = true;
    if (!sourceFolderPath || !targetFolderPath || !watermarkPath) {
      alert("Please select all the required fields");
      return;
    }
    ProcessImages(
      watermarkPath,
      sourceFolderPath,
      targetFolderPath,
      watermarkPosition,
      watermarkOpacity / 100,
      watermarkScale / 100
    ).then((result) => {
      showLoader = false;
      if (result.status && result.status == "error") {
        alert(result.message);
      } else {
        alert(result.message);
      }
    });
  }

  function mainButtonFunction() {
    if (sourceFolderPath == "" || targetFolderPath == "" || watermarkPath == "") {
      alert("Please select all the required fields");
    } else if (numberOfSourceFiles == 0) {
      alert("No images found in source folder!");
    } else if (changedSettings) {
      generatePreview();
    } else {
      processFiles();
    }
  }
</script>

<main>
  {#if showLoader}
    <div class="loader">
      <Loader />
      <div class="loadingtext">Processing image</div>
    </div>
  {/if}
  <Header {title} {version} {author} {authorLink} />
  <div class="maincontainer">
    <div class="settingscol">
      <div class="pathconfig">
        <Card backgroundcolor={"#ff0000"}>
          <Pathconfig
            bind:sourceFolderPath
            bind:targetFolderPath
            bind:watermarkPath
            bind:numberOfSourceFiles
            bind:numberOfTargetFiles
            bind:watermarkpreviewImage
          />
        </Card>
      </div>
      <div class="settings">
        <Card backgroundcolor={"#0000ff"}>
          <Settings bind:watermarkOpacity bind:watermarkScale bind:watermarkPosition bind:changedSettings />
        </Card>
      </div>
    </div>

    <div class="mainview">
      <img bind:this={watermarkedpreviewImage} src="" alt="Preview of the result" />
      <Button bind:text={mainbuttontext} callback={mainButtonFunction} />
    </div>
  </div>
</main>

<style>
  main {
    background: white;
    height: 100%;
    width: 100%;
  }
  .loader {
    position: absolute;
    width: 100%;
    height: 100%;
    background-color: black;
    z-index: 100;
    opacity: 0.5;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .loadingtext {
    font-size: 40px;
    padding-top: 20px;
  }

  .mainview {
    width: 75%;
    height: 100%;
  }
  .mainview img {
    max-width: 95%;
    border-radius: 20px;
    margin: 10px;
    box-shadow: 0px 0px 20px #4d5056;
    /* width: 733.875px; */
    height: 489.25px;
  }

  .settingscol {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    width: 25%;
  }
  .maincontainer {
    display: flex;
    flex-direction: row;
    position: absolute;
    top: 30px;
    bottom: 0;
    width: 100%;
  }
</style>
