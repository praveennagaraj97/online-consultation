import { Component, Input } from '@angular/core';
import { BreadcrumbPath } from 'src/app/types/app.types';

@Component({
  selector: 'app-breadcrumb-component',
  template: `
    <div class="flex items-center" [ngClass]="className">
      <a
        routerLink="/dashboard"
        title="dashboard"
        class="hover:scale-105 smooth-animate"
      >
        <app-home-icon
          className="w-6 h-6 dark:fill-gray-50 dark:stroke-gray-50"
        ></app-home-icon>
      </a>

      <ng-container *ngFor="let path of paths">
        <app-chevron-icon
          className="w-6 h-6 dark:fill-gray-50 dark:stroke-gray-50 mx-1"
        ></app-chevron-icon>
        <a
          [routerLink]="path.path"
          [title]="path.title"
          class="text-sm anchor-link hover:scale-105"
        >
          {{ path.title }}
        </a>
      </ng-container>
    </div>
  `,
})
export class BreadcrumbComponent {
  @Input() paths: BreadcrumbPath[] = [];
  @Input() className = '';
}
