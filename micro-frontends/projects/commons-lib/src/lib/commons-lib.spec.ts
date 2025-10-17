import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CommonsLib } from './commons-lib';

describe('CommonsLib', () => {
  let component: CommonsLib;
  let fixture: ComponentFixture<CommonsLib>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CommonsLib]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CommonsLib);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
