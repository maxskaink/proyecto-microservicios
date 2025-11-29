import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IsLoading } from './is-loading';

describe('IsLoading', () => {
  let component: IsLoading;
  let fixture: ComponentFixture<IsLoading>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IsLoading]
    })
    .compileComponents();

    fixture = TestBed.createComponent(IsLoading);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
