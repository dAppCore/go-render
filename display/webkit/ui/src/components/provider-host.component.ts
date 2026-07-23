// SPDX-Licence-Identifier: EUPL-1.2

import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  ElementRef,
  Input,
  OnChanges,
  OnInit,
  Renderer2,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ApiConfigService } from '../services/api-config.service';
import { ProviderDiscoveryService } from '../services/provider-discovery.service';

/**
 * ProviderHostComponent renders any custom element by tag name using
 * Angular's Renderer2 for safe DOM manipulation. It reads the :provider
 * route parameter to look up the element tag from the discovery service.
 */
@Component({
  selector: 'provider-host',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  template: '<div #container class="provider-host"></div>',
  styles: [
    `
      :host {
        display: block;
        width: 100%;
        height: 100%;
      }
      .provider-host {
        width: 100%;
        height: 100%;
      }
    `,
  ],
})
export class ProviderHostComponent implements OnInit, OnChanges {
  /** The custom element tag to render. Can be set via input or route param. */
  @Input() tag = '';

  /** API URL attribute passed to the custom element. */
  @Input() apiUrl = '';

  @ViewChild('container', { static: true }) container!: ElementRef;

  constructor(
    private renderer: Renderer2,
    private route: ActivatedRoute,
    private apiConfig: ApiConfigService,
    private providerService: ProviderDiscoveryService,
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      const providerName = params['provider'];
      if (providerName) {
        const provider = this.providerService
          .providers()
          .find((p) => p.name.toLowerCase() === providerName.toLowerCase());
        if (provider?.element?.tag) {
          this.tag = provider.element.tag;
          this.renderElement();
        }
      }
    });
  }

  ngOnChanges(): void {
    this.renderElement();
  }

  private renderElement(): void {
    if (!this.tag || !this.container) {
      return;
    }

    const native = this.container.nativeElement;

    // Clear previous element safely
    while (native.firstChild) {
      this.renderer.removeChild(native, native.firstChild);
    }

    // Create and append the custom element
    const el = this.renderer.createElement(this.tag);
    const url = this.apiUrl || this.apiConfig.baseUrl;
    if (url) {
      this.renderer.setAttribute(el, 'api-url', url);
    }
    this.renderer.appendChild(native, el);
  }
}
