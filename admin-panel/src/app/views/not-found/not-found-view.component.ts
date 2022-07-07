import { Component } from '@angular/core';

@Component({
  selector: 'app-page-not-found-view',
  template: `
    <section
      class="flex items-center h-screen w-screen dark:bg-gray-900 bg-indigo-100 dark:text-gray-100"
    >
      <div
        class="container flex flex-col items-center justify-center px-5 mx-auto my-8"
      >
        <div class="max-w-md text-center">
          <h2 class="mb-8 font-extrabold text-9xl dark:text-gray-600">
            <span class="sr-only">Error</span>404
          </h2>
          <p class="text-2xl font-semibold md:text-3xl">
            Sorry, we couldn't find this page.
          </p>
          <p class="mt-4 mb-8 dark:text-gray-400">
            But dont worry, you can find plenty of other things on our homepage.
          </p>
          <a
            rel="noopener noreferrer"
            href="#"
            class="px-8 font-semibold action-btn  flex items-center justify-center shadow-md text-sm rounded-md py-2 mt-1 w-fit mx-auto block"
            >Back to homepage</a
          >
        </div>
      </div>
    </section>
  `,
})
export class PageNotFoundViewComponent {}
