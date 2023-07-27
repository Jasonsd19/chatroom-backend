<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import Swal from "sweetalert2";

  interface UserMessage {
    username: string;
    message: string;
  }

  let username: string = "";
  let sendText: HTMLTextAreaElement;
  let messages: UserMessage[] = [];

  let connected = false;
  let ws: WebSocket | null = null;

  const onPressEnter = (ev: KeyboardEvent) => {
    if (ev.code === "Enter" && !ev.shiftKey && ws) {
      ev.preventDefault();
      const m: UserMessage = { username, message: sendText.value };
      if (m && m.message.length <= 256) {
        ws.send(JSON.stringify(m));
        sendText.value = "";
      }
    }
  };

  const onChangeUsername = (e: Event) => {
    const event = e.target as HTMLInputElement;
    username = event.value;
  };

  onMount(() => {
    window.addEventListener("keypress", onPressEnter);
    return () => window.removeEventListener("keypress", onPressEnter);
  });

  onDestroy(() => {
    if (ws && ws.readyState === ws.OPEN) ws.close();
  });

  const connect = () => {
    if (username.length <= 8) {
      Swal.fire({
        title: "Error",
        text: "Username must be at least 8 characters long.",
        icon: "error",
        confirmButtonText: "Close",
      });
      return;
    }

    if ("WebSocket" in window || "MozWebSocket" in window) {
      ws = new WebSocket("wss://jasonsd-chatroom-backend-f0a0780d9803.herokuapp.com/ws");
      ws.onopen = () => {
        connected = true;
      };

      ws.onclose = () => {
        connected = false;
      };

      ws.onerror = () => {
        connected = false;
        Swal.fire({
          title: "Error",
          text: "We ran into an unexpected error, please refresh the page and try again.",
          icon: "error",
          confirmButtonText: "Close",
        });
      };

      ws.onmessage = (event) => {
        const message = JSON.parse(event.data) as UserMessage;
        if (messages.length >= 50) messages = [message, ...messages.slice(0, -1)];
        else messages = [message, ...messages];
      };
    } else {
      Swal.fire({
        title: "Error",
        text: "Your browser does not support websockets!",
        icon: "error",
        confirmButtonText: "Close",
      });
    }
  };
</script>

<div class="mainContainer">
  {#if connected}
    <div class="chatroomContainer">
      <div class="messageContainer">
        {#if messages.length}
          {#each messages as message, i (i)}
            <div class="message">
              {`${message.username}: ${message.message}`}
            </div>
          {/each}
        {:else}
          <div class="message">
            {"No messages yet!"}
          </div>
        {/if}
      </div>
      <div class="textContainer">
        <textarea maxlength={200} bind:this={sendText} />
      </div>
    </div>
  {:else}
    <div class="usernameContainer">
      <label for="username">Username: </label>
      <input inputmode="text" on:input={onChangeUsername} />
    </div>
    <div class="joinButtonContainer">
      <button on:click={connect}>Join Chatroom!</button>
    </div>
  {/if}
</div>

<style>
  .mainContainer {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    overflow-x: hidden;
    flex-wrap: wrap;
    align-content: center;
    justify-content: center;
  }

  .chatroomContainer {
    display: flex;
    height: 60vh;
    width: 40vw;
    flex-wrap: wrap;
    align-content: center;
    justify-content: center;
  }

  .messageContainer {
    display: flex;
    flex-direction: column-reverse;
    height: 100%;
    width: 100%;
    overflow-x: hidden;
    overflow-y: visible;
  }

  .message {
    display: flex;
    flex-wrap: wrap;
    width: 100%;
    justify-content: left;
    padding: 1vh;
  }

  .textContainer {
    display: flex;
    flex-wrap: wrap;
    align-content: center;
    justify-content: center;
    height: 10%;
    width: 50%;
    padding: 1vw;
  }

  .textContainer textarea {
    height: 100%;
    width: 100%;
    resize: none;
  }

  .usernameContainer {
    width: 20%;
    height: 3%;
    display: flex;
    font-size: 1.5vw;
    padding: 5vw;
    justify-content: space-evenly;
    align-items: center;
  }

  .usernameContainer label {
    margin-bottom: 0.5vw;
    margin-right: 0.5vw;
  }

  .usernameContainer input {
    height: 100%;
    width: 100%;
    font-size: 1.5vw;
  }

  .joinButtonContainer {
    display: flex;
    flex-wrap: wrap;
    align-content: center;
    justify-content: center;
  }
</style>
