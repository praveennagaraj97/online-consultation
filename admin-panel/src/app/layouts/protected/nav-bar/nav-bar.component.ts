import {
  animate,
  state,
  style,
  transition,
  trigger,
} from '@angular/animations';
import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-nav-bar',
  templateUrl: 'nav-bar.component.html',
  animations: [
    trigger('expand', [
      state('open', style({ width: '16rem' })),
      state('closed', style({ width: '6rem' })),
      transition('open => closed', animate('0.2s')),
      transition('closed => open', animate('0.3s')),
    ]),
  ],
})
export class NavBarComponent {
  expandMenu = false;
  expandState: 'open' | 'closed' = 'closed';

  @Output() onExpand: EventEmitter<boolean> = new EventEmitter(false);

  toggleMenu(state: boolean) {
    this.expandState = state ? 'open' : 'closed';
    this.expandMenu = state;
    this.onExpand.emit(state);
  }
}
