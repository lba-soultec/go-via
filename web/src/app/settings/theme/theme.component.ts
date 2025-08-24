import { Component } from '@angular/core';
import { ApiService } from '../../api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-theme',
  templateUrl: './theme.component.html',
  styleUrls: ['./theme.component.scss'],
  standalone: false
})
export class ThemeComponent {
  selectedFile: File | null = null;
  previewUrl: string | ArrayBuffer | null = null;
  successMessage = '';
  errorMessage = '';

  constructor(private api: ApiService, private router: Router) {}

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      if (!['image/png', 'image/jpeg'].includes(file.type)) {
        this.errorMessage = 'Only JPG and PNG files are allowed.';
        this.selectedFile = null;
        this.previewUrl = null;
        return;
      }
      if (file.size > 5 * 1024 * 1024) {
        this.errorMessage = 'File size must not exceed 5MB.';
        this.selectedFile = null;
        this.previewUrl = null;
        return;
      }
      this.selectedFile = file;
      this.errorMessage = '';
      const reader = new FileReader();
      reader.onload = e => this.previewUrl = reader.result;
      reader.readAsDataURL(file);
    }
  }

  onSubmit(): void {
    if (!this.selectedFile) return;
    const formData = new FormData();
    formData.append('background', this.selectedFile);
    this.api.uploadBackgroundImage(formData).subscribe({
      next: () => {
        this.successMessage = 'Background image updated successfully! The new background will be visible on next page load.';
        this.errorMessage = '';
        // Navigate away and back to refresh the background
        setTimeout(() => {
          this.router.navigateByUrl('/', { skipLocationChange: true }).then(() => {
            this.router.navigate(['/settings/theme']);
          });
        }, 2000);
      },
      error: err => {
        this.errorMessage = 'Failed to upload image. Please try again.';
        this.successMessage = '';
      }
    });
  }
}
