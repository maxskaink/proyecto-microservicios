import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-arrow-left',
  imports: [CommonModule],
  templateUrl: './arrow-left.html',
  styleUrl: './arrow-left.css',
})
export class ArrowLeft {
  @Input() url!: string;
  constructor(private router: Router){}
  goBack(): void {
    if (this.url) {
      this.router.navigate([this.url]); 
    } else {
      window.history.back(); 
    }
  }
}
