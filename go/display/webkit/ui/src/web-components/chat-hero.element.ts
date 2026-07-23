const template = document.createElement('template');
template.innerHTML = `
  <style>
    :host {
      display: block;
      border: 1px solid rgba(251, 191, 36, 0.16);
      border-radius: 24px;
      padding: 1.2rem 1.35rem;
      background:
        linear-gradient(135deg, rgba(120, 53, 15, 0.38), rgba(15, 23, 42, 0.88)),
        radial-gradient(circle at top right, rgba(125, 211, 252, 0.18), transparent 36%);
      box-shadow: 0 18px 50px rgba(2, 6, 23, 0.32);
    }
    .eyebrow {
      margin: 0;
      color: #fbbf24;
      text-transform: uppercase;
      letter-spacing: 0.18em;
      font: 600 0.72rem/1.2 "Avenir Next Condensed", "Gill Sans", sans-serif;
    }
    h1 {
      margin: 0.35rem 0 0;
      color: #f8fafc;
      font: 700 clamp(1.9rem, 3vw, 3rem)/1 "Iowan Old Style", "Palatino Linotype", serif;
    }
    .subtitle {
      margin: 0.55rem 0 0;
      color: #cbd5e1;
      font: 500 0.95rem/1.5 "Avenir Next", "Segoe UI", sans-serif;
    }
  </style>
  <p class="eyebrow"></p>
  <h1></h1>
  <p class="subtitle"></p>
`;

class CoreChatHeroElement extends HTMLElement {
  static get observedAttributes(): string[] {
    return ['eyebrow', 'title', 'subtitle'];
  }

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: 'open' });
    shadowRoot.appendChild(template.content.cloneNode(true));
  }

  connectedCallback(): void {
    this.render();
  }

  attributeChangedCallback(): void {
    this.render();
  }

  private render(): void {
    const root = this.shadowRoot;
    if (!root) {
      return;
    }
    const eyebrow = root.querySelector('.eyebrow');
    const title = root.querySelector('h1');
    const subtitle = root.querySelector('.subtitle');
    if (eyebrow) eyebrow.textContent = this.getAttribute('eyebrow') || 'CoreGUI Chat';
    if (title) title.textContent = this.getAttribute('title') || 'Local chat';
    if (subtitle)
      subtitle.textContent =
        this.getAttribute('subtitle') || 'Local-first chat, UI events, and sidecar tooling';
  }
}

export function registerChatHeroElement(): void {
  if (!customElements.get('core-chat-hero')) {
    customElements.define('core-chat-hero', CoreChatHeroElement);
  }
}
