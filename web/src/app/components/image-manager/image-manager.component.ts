import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-image-manager',
  templateUrl: './image-manager.component.html',
  styleUrls: ['./image-manager.component.scss'],
  standalone: false
})
export class ImageManagerComponent {
  @Input() images: any[] = [];
  @Output() imagesChange = new EventEmitter<any[]>();

  hash = '';
  description = '';
  selectedFiles: FileList | null = null;
  progress = 0;
  message = '';

  selectFile(event: Event) {
    const input = event.target as HTMLInputElement;
    this.selectedFiles = input.files;
  }

  upload() {
    if (!this.selectedFiles || this.selectedFiles.length === 0) {
      this.message = 'No file selected!';
      return;
    }
    // Simulate upload logic
    const file = this.selectedFiles[0];
    const newImage = {
      id: Date.now(),
      iso_image: file.name,
      size: Math.round(file.size / (1024 * 1024)),
      description: this.description,
      hash: this.hash
    };
    this.images.push(newImage);
    this.imagesChange.emit(this.images);
    this.progress = 100;
    this.message = 'Upload successful!';
    // Reset form
    this.selectedFiles = null;
    this.hash = '';
    this.description = '';
    setTimeout(() => {
      this.progress = 0;
      this.message = '';
    }, 2000);
  }

  remove(id: number) {
    this.images = this.images.filter(img => img.id !== id);
    this.imagesChange.emit(this.images);
  }
}
