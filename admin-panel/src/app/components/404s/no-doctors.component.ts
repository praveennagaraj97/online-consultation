import { Component } from '@angular/core';

@Component({
  selector: 'app-no-doctors-component',
  template: `
    <div class="py-4 flex flex-col items-center justify-center">
      <img src="/assets/doctor-grp.png" alt="no-doctors" />
      <h2 class="mt-2 font-medium">Couldn't find any doctors.</h2>
    </div>
  `,
})
export class NoDoctorsFoundComponent {}
