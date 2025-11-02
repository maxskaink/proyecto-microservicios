import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ArrowLeft } from './arrow-left';

describe('ArrowLeft', () => {
  let component: ArrowLeft;
  let fixture: ComponentFixture<ArrowLeft>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ArrowLeft]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ArrowLeft);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
