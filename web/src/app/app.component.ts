import { Component, OnInit, Renderer2 } from '@angular/core';
import { ApiService } from './api.service';
import { Router, NavigationEnd } from '@angular/router';
import { AuthService } from './auth.service';
import { filter } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  standalone: false
})
export class AppComponent implements OnInit {
  version: any;
  username: string;
  showAbout: boolean = false;

  constructor(
    private apiService: ApiService, 
    public router: Router, 
    private authService: AuthService,
    private renderer: Renderer2
  ) {}

  ngOnInit(): void {

    this.authService.username$.subscribe(username => {
      this.username = username;
    });

    this.apiService.getVersion().subscribe((data: any) => {
      this.version = data;
      console.log(this.version);
    });

    if (!this.authService.isLoggedIn()) {
      this.username = null;
      this.router.navigate(['/login']);
    } else {
      // Load background image for authenticated users
      this.loadBackgroundImage();
    }

    // Listen for route changes to apply background when user logs in
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe(() => {
      if (this.authService.isLoggedIn() && this.router.url !== '/login') {
        this.loadBackgroundImage();
      }
    });
  }

  private loadBackgroundImage(): void {
    this.apiService.getThemeImage().subscribe({
      next: (blob) => {
        if (blob && blob.size > 0) {
          const reader = new FileReader();
          reader.onload = () => {
            this.setBackground(reader.result as string);
          };
          reader.readAsDataURL(blob);
        } else {
          this.setBackground('/assets/background.png');
        }
      },
      error: () => {
        this.setBackground('/assets/background.png');
      }
    });
  }

  private setBackground(url: string): void {
    const el = document.querySelector('.main-container') as HTMLElement;
    if (el) {
      this.renderer.setStyle(el, 'background-image', `url('${url}')`);
      this.renderer.setStyle(el, 'background-size', 'cover');
      this.renderer.setStyle(el, 'background-position', 'center');
      this.renderer.setStyle(el, 'background-repeat', 'no-repeat');
    }
  }

  logout() {
    this.authService.logout();
    this.username = null;
  }

  showAboutModel(mode, id=null) {
    this.showAbout = mode;

  }
}


