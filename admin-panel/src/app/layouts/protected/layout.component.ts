import { Component } from '@angular/core';

@Component({
  selector: 'app-protected-layout',
  templateUrl: 'layout.component.html',
})
export class ProtectedLayoutComponent {
  menuExpanded = false;

  onMenuExpand(state: boolean) {
    this.menuExpanded = state;
  }
}
